
# Qoin Go Helper


[![logo](https://avatars.githubusercontent.com/u/72009988?s=200&v=4)](https://www.qoin.id/)


# Introduction
This repo for Helper Golang services.

# ⚙️ Installation

### via go modules
```git
go get github.com/Qoin-Digital-Indonesia/qoingohelper
```

## 🎯 Features

-   Logger [logger](https://github.com/Qoin-Digital-Indonesia/qoingohelper/blob/master/logger.go)
-   Data [data](https://github.com/Qoin-Digital-Indonesia/qoingohelper/blob/master/data.go)
-   Middleware [middleware](https://github.com/Qoin-Digital-Indonesia/qoingohelper/blob/master/middleware.go)
-   Response V1 snake_case
-   Response V2 CamlCase
-   RSA PEM keys + RSA-SHA256 signatures [qoincrypto](qoincrypto/rsa.go) (replaces `github.com/supriadi-yusuf/adicrypto`)

# ⚡️ Usage
## 📖 Logger
```go
import (
     ...
     github.com/Qoin-Digital-Indonesia/qoingohelper
     ...
)

func main(){
     ...
      qoingohelper.LoggerInfo("This is info")
	 qoingohelper.LoggerSuccess("This is success")
	 qoingohelper.LoggerError(errors.New("this is info")) or qoingohelper.LoggerError(err)
	 qoingohelper.LoggerWarning("This is warning")     ​...

}

```
## 📖 Logger Sentry
```go
import (
     ...
     github.com/Qoin-Digital-Indonesia/qoingohelper
     ...
)

func main(){
     ...
      qoingohelper.InitSentry("this is url sentry", "environment (develop, staging , production)", "tags release (1.0.0 , 1.0.12)", debug (true,false))
     ...

     ...
     // Send Error
     if err != nil {
		qoingohelper.SendSentryError(err, "saas-be-research-manager", "User", "divide")
          // 1. Error
          // 2. Service
          // 3. Module
          // 4. Function
	}
     ...

     // Send Message
     ...
     qoingohelper.SendSentryMessage("this is sample message", "saas-be-research-manager", "User", "divide")
          // 1. Message
          // 2. Service
          // 3. Module
          // 4. Function
     ...

     // Send Event
     ...
     qoingohelper.SendSentryEvent(*sentry.Event)
          // 1. Sentry Event
          // 2. Service
          // 3. Module
          // 4. Function
     ...
}

```


###
* Result

```diff
- logger error
+ logger warning
! logger success
# logger info
```

## 📖 Data

## 📖 qoincrypto (RSA PEM + RSA-SHA256 signature)
Replaces `github.com/supriadi-yusuf/adicrypto`. Same signature format (RSASSA-PKCS1-v1_5 over SHA-256),
so signatures made by either library verify with the other. Standard library only; never logs key material.

```go
import "github.com/Qoin-Digital-Indonesia/qoingohelper/v2/qoincrypto"

priv, err := qoincrypto.ParsePrivateKeyPEM([]byte(os.Getenv("APP_PRIVATE_KEY"))) // PKCS#1 or PKCS#8
pub, err := qoincrypto.ParsePublicKeyPEM([]byte(os.Getenv("APP_PUBLIC_KEY")))    // PKIX or PKCS#1
if qoincrypto.PublicKeyBits(pub) < qoincrypto.RecommendedRSAKeyBits { /* warn: rotate to >= 2048 bits */ }

sig, err := qoincrypto.SignSHA256(priv, payload)
err = qoincrypto.VerifySHA256(pub, payload, sig) // nil = valid, qoincrypto.ErrInvalidSignature = mismatch
```

| adicrypto | qoincrypto |
|---|---|
| `PEMBytesToPrivateKey` | `ParsePrivateKeyPEM` (encrypted PEM rejected) |
| `PEMBytesToPublicKey` | `ParsePublicKeyPEM` |
| hasher + `CreateSignCreator(...).GenerateSignature` | `SignSHA256` (nil key → error, not panic) |
| hasher + `CreateSignVerificator(...).VerifySignature` | `VerifySHA256` (mismatch → `ErrInvalidSignature`; nil/weak key → other error) |
