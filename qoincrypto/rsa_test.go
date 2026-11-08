package qoincrypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"sync"
	"testing"
)

var (
	testKeyOnce sync.Once
	testKey     *rsa.PrivateKey
)

// kunci dibuat saat test berjalan; tidak ada materi kunci di repo.
func key(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	testKeyOnce.Do(func() {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		testKey = k
	})
	return testKey
}

func encode(typ string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der})
}

func TestParsePrivateKeyPEM(t *testing.T) {
	k := key(t)
	pkcs8, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{
		"pkcs1":             encode("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(k)),
		"pkcs8":             encode("PRIVATE KEY", pkcs8),
		"pkcs1 wrong label": encode("PRIVATE KEY", x509.MarshalPKCS1PrivateKey(k)),
		"surrounding text":  append(append([]byte("\n  "), encode("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(k))...), '\n'),
	}
	for name, in := range cases {
		got, err := ParsePrivateKeyPEM(in)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", name, err)
			continue
		}
		if !got.Equal(k) {
			t.Errorf("%s: parsed key differs", name)
		}
	}
}

func TestParsePrivateKeyPEMErrors(t *testing.T) {
	ec, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecDER, _ := x509.MarshalPKCS8PrivateKey(ec)
	encrypted := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: map[string]string{"Proc-Type": "4,ENCRYPTED", "DEK-Info": "AES-256-CBC,00"},
		Bytes:   []byte{1, 2, 3},
	})
	cases := []struct {
		name string
		in   []byte
		want error
	}{
		{"empty", nil, ErrEmptyKey},
		{"not pem", []byte("hello"), ErrNoPEMBlock},
		{"legacy encrypted", encrypted, ErrEncryptedPEM},
		{"pkcs8 encrypted", encode("ENCRYPTED PRIVATE KEY", []byte{1}), ErrEncryptedPEM},
		{"ec key", encode("PRIVATE KEY", ecDER), ErrNotRSAKey},
		{"garbage der", encode("RSA PRIVATE KEY", []byte{1, 2, 3}), nil},
	}
	for _, c := range cases {
		got, err := ParsePrivateKeyPEM(c.in)
		if err == nil || got != nil {
			t.Errorf("%s: expected error, got key=%v err=%v", c.name, got != nil, err)
			continue
		}
		if c.want != nil && !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestParsePublicKeyPEM(t *testing.T) {
	pub := &key(t).PublicKey
	pkix, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{
		"pkix":             encode("PUBLIC KEY", pkix),
		"pkcs1":            encode("RSA PUBLIC KEY", x509.MarshalPKCS1PublicKey(pub)),
		"pkix wrong label": encode("RSA PUBLIC KEY", pkix), // format yang dihasilkan adicrypto.SavePublicKey
	}
	for name, in := range cases {
		got, err := ParsePublicKeyPEM(in)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", name, err)
			continue
		}
		if !got.Equal(pub) {
			t.Errorf("%s: parsed key differs", name)
		}
	}

	ec, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecDER, _ := x509.MarshalPKIXPublicKey(&ec.PublicKey)
	if _, err := ParsePublicKeyPEM(encode("PUBLIC KEY", ecDER)); !errors.Is(err, ErrNotRSAKey) {
		t.Errorf("ec public key: err = %v, want ErrNotRSAKey", err)
	}
	if _, err := ParsePublicKeyPEM(nil); !errors.Is(err, ErrEmptyKey) {
		t.Errorf("empty: err = %v, want ErrEmptyKey", err)
	}
	if _, err := ParsePublicKeyPEM(encode("PUBLIC KEY", []byte{1, 2})); err == nil {
		t.Error("garbage der: expected error")
	}
}

func TestSignVerifySHA256(t *testing.T) {
	k := key(t)
	payload := []byte(`{"id":"123","amount":"10000"}`)

	sig, err := SignSHA256(k, payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Kompatibel dengan format adicrypto (PKCS#1 v1.5 + SHA-256, deterministik).
	h := sha256.Sum256(payload)
	want, _ := rsa.SignPKCS1v15(nil, k, crypto.SHA256, h[:])
	if string(sig) != string(want) {
		t.Fatal("signature differs from RSASSA-PKCS1-v1_5/SHA-256 reference")
	}

	if err := VerifySHA256(&k.PublicKey, payload, sig); err != nil {
		t.Fatalf("verify valid signature: %v", err)
	}

	tampered := append([]byte{}, payload...)
	tampered[0] ^= 1
	if err := VerifySHA256(&k.PublicKey, tampered, sig); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("tampered payload: err = %v, want ErrInvalidSignature", err)
	}
	if err := VerifySHA256(&k.PublicKey, payload, sig[:len(sig)-1]); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("truncated signature: err = %v, want ErrInvalidSignature", err)
	}
	if err := VerifySHA256(&k.PublicKey, payload, nil); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("empty signature: err = %v, want ErrInvalidSignature", err)
	}
}

func TestNilAndWeakKeys(t *testing.T) {
	if _, err := SignSHA256(nil, []byte("x")); !errors.Is(err, ErrNilKey) {
		t.Errorf("sign nil key: err = %v", err)
	}
	if err := VerifySHA256(nil, []byte("x"), []byte("y")); !errors.Is(err, ErrNilKey) {
		t.Errorf("verify nil key: err = %v", err)
	}
	weak := &rsa.PublicKey{N: new(big.Int).Lsh(big.NewInt(1), 511), E: 65537}
	err := VerifySHA256(weak, []byte("x"), []byte("y"))
	if err == nil || errors.Is(err, ErrInvalidSignature) {
		t.Errorf("512-bit key: err = %v, want configuration error", err)
	}
	if PublicKeyBits(weak) != 512 || PublicKeyBits(nil) != 0 {
		t.Error("PublicKeyBits")
	}
}
