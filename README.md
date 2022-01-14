
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


###
* Result

```diff
- logger error
+ logger warning
! logger success
# logger info
```

## 📖 Data
