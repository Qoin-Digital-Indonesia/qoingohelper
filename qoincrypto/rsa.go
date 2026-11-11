// Package qoincrypto berisi helper kriptografi RSA yang dipakai service QoinPayout:
// parsing kunci PEM dan tanda tangan RSA PKCS#1 v1.5 dengan SHA-256.
//
// Package ini menggantikan github.com/supriadi-yusuf/adicrypto. Format tanda tangan sama
// persis (RSASSA-PKCS1-v1_5 + SHA-256), jadi tanda tangan lama tetap valid dan sebaliknya.
// Perbedaannya:
//   - PEM terenkripsi (legacy "Proc-Type: 4,ENCRYPTED", RFC 1423) ditolak; adicrypto mencoba
//     mendekripsinya dengan password kosong memakai API yang sudah deprecated dan tidak aman.
//   - Kunci privat PKCS#1 dan PKCS#8, kunci publik PKIX dan PKCS#1 diterima.
//   - Kunci nil, input kosong, atau kunci non-RSA dikembalikan sebagai error, bukan panic.
//   - Verifikasi membedakan "tanda tangan tidak cocok" (ErrInvalidSignature) dari error lain
//     (kunci tidak ada, kunci terlalu lemah untuk runtime Go).
//
// Package ini hanya memakai standard library dan tidak pernah mencetak materi kunci.
package qoincrypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// RecommendedRSAKeyBits adalah ukuran kunci minimum yang direkomendasikan (NIST SP 800-57).
const RecommendedRSAKeyBits = 2048

var (
	// ErrEmptyKey: input PEM kosong (biasanya env belum di-set).
	ErrEmptyKey = errors.New("qoincrypto: key is empty")
	// ErrNoPEMBlock: input tidak berisi blok PEM.
	ErrNoPEMBlock = errors.New("qoincrypto: no PEM block found")
	// ErrEncryptedPEM: blok PEM terenkripsi tidak didukung.
	ErrEncryptedPEM = errors.New("qoincrypto: encrypted PEM keys are not supported")
	// ErrNotRSAKey: kunci valid tetapi bukan RSA.
	ErrNotRSAKey = errors.New("qoincrypto: key is not an RSA key")
	// ErrNilKey: kunci nil diberikan ke Sign/Verify.
	ErrNilKey = errors.New("qoincrypto: key is nil")
	// ErrInvalidSignature: tanda tangan tidak cocok dengan payload dan kunci publik.
	ErrInvalidSignature = errors.New("qoincrypto: invalid signature")
)

func decodePEM(pemBytes []byte) (*pem.Block, error) {
	// Nilai dari env sering diawali spasi/baris kosong; pem.Decode menolak "-----BEGIN" yang menjorok.
	pemBytes = bytes.TrimSpace(pemBytes)
	if len(pemBytes) == 0 {
		return nil, ErrEmptyKey
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrNoPEMBlock
	}
	if block.Type == "ENCRYPTED PRIVATE KEY" || block.Headers["Proc-Type"] != "" || block.Headers["DEK-Info"] != "" {
		return nil, ErrEncryptedPEM
	}
	return block, nil
}

// ParsePrivateKeyPEM mengubah kunci privat RSA berformat PEM (PKCS#1 "RSA PRIVATE KEY"
// atau PKCS#8 "PRIVATE KEY") menjadi *rsa.PrivateKey. Kunci divalidasi (Precompute/Validate
// dijalankan oleh parser x509).
func ParsePrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, err := decodePEM(pemBytes)
	if err != nil {
		return nil, err
	}

	// Coba berdasarkan isi, bukan label: kunci lama kadang diberi label yang tidak sesuai.
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("qoincrypto: parse private key (%s): not a valid PKCS#1 or PKCS#8 key", block.Type)
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotRSAKey
	}
	return key, nil
}

// ParsePublicKeyPEM mengubah kunci publik RSA berformat PEM (PKIX "PUBLIC KEY" atau
// PKCS#1 "RSA PUBLIC KEY") menjadi *rsa.PublicKey.
func ParsePublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, err := decodePEM(pemBytes)
	if err != nil {
		return nil, err
	}

	if parsed, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		key, ok := parsed.(*rsa.PublicKey)
		if !ok {
			return nil, ErrNotRSAKey
		}
		return key, nil
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("qoincrypto: parse public key (%s): not a valid PKIX or PKCS#1 key", block.Type)
	}
	return key, nil
}

// PublicKeyBits mengembalikan ukuran modulus kunci publik dalam bit (0 bila nil).
// Dipakai service untuk memberi peringatan saat start bila kunci < RecommendedRSAKeyBits.
func PublicKeyBits(key *rsa.PublicKey) int {
	if key == nil || key.N == nil {
		return 0
	}
	return key.N.BitLen()
}

// SignSHA256 membuat tanda tangan RSASSA-PKCS1-v1_5 atas SHA-256(payload).
func SignSHA256(privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, ErrNilKey
	}
	hashed := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("qoincrypto: sign: %w", err)
	}
	return signature, nil
}

// VerifySHA256 memverifikasi tanda tangan RSASSA-PKCS1-v1_5 atas SHA-256(payload).
// Mengembalikan nil bila valid, ErrInvalidSignature bila tidak cocok, atau error lain bila
// verifikasi tidak bisa dilakukan (kunci nil, kunci terlalu kecil untuk runtime Go).
func VerifySHA256(publicKey *rsa.PublicKey, payload, signature []byte) error {
	if publicKey == nil || publicKey.N == nil {
		return ErrNilKey
	}
	// Runtime Go (>= 1.24) menolak kunci < 1024 bit; laporkan sebagai error konfigurasi,
	// bukan sebagai tanda tangan salah, supaya mudah didiagnosis.
	if bits := PublicKeyBits(publicKey); bits < 1024 {
		return fmt.Errorf("qoincrypto: verify: public key is %d bits, below the 1024-bit minimum enforced by Go", bits)
	}
	hashed := sha256.Sum256(payload)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return ErrInvalidSignature
	}
	return nil
}
