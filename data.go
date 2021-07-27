package qoingohelper

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Validates *validator.Validate

func DateParse(d interface{}, format string) *time.Time {

	if d == nil {
		return nil
	}

	dates, err := time.Parse(format, fmt.Sprintf("%v", d))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &dates
}

func DateNil(dt interface{}) interface{} {

	switch v := dt.(type) {
	case string, time.Time:
		dates, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", fmt.Sprintf("%v", v))
		if err != nil || dates.IsZero() {
			dates, err = time.Parse("2006-01-02", fmt.Sprintf("%v", v))
			if err != nil || dates.IsZero() {
				dates, err = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v", v))
				if err != nil || dates.IsZero() {
					dates, err = time.Parse("2006-01-02T00:00:00Z", fmt.Sprintf("%v", v))
					if err != nil || dates.IsZero() {
						return nil
					}
					return dates.Format("2006-01-02")
				}
				return dates.Format("2006-01-02 15:04:05")
			}
			return dates.Format("2006-01-02")
		}
		return dates.Format("2006-01-02 15:04:05")
	default:
		return nil
	}
}

func Validate(c echo.Context, i interface{}) (interface{}, error) {
	err := c.Validate(i)
	if err != nil {
		messageError := make(map[string]interface{})
		for _, err := range err.(validator.ValidationErrors) {
			messageError[err.StructField()] = err.StructField() + " Is " + err.ActualTag()
		}
		return messageError, err
	}
	return nil, err
}

func ValidateRPC(i interface{}) (interface{}, error) {
	Validates = validator.New()

	err := Validates.Struct(i)
	if err != nil {
		messageError := make(map[string]interface{})
		for _, err := range err.(validator.ValidationErrors) {
			messageError[err.StructField()] = err.StructField() + " Is " + err.ActualTag()
		}
		return messageError, err
	}
	return nil, err
}

func Bind(form interface{}, structs interface{}) interface{} {
	v, err := json.Marshal(form)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(v, &structs)

	if err != nil {
		return nil
	}

	return structs
}

func UUID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	uuid := fmt.Sprintf("%x", b[0:10])
	return uuid
}

func RandStringBytes(n int) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mrand.Intn(len(letterBytes))]
	}
	return string(b)
}

func InArray(val interface{}, arrays interface{}) bool {
	kind := reflect.TypeOf(arrays).Kind()
	values := reflect.ValueOf(arrays)

	if kind == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if fmt.Sprint(val) == fmt.Sprint(values.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

func JsonEncode(data interface{}) string {
	val, err := json.Marshal(data)
	if err != nil {
		LoggerError(err)
		return ""
	}

	return string(val)
}

func JsonDecode(data interface{}) (maps map[string]interface{}, err error) {
	toMapping := make(map[string]interface{})
	jsonString, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonString, &toMapping)

	if err != nil {
		return nil, err
	}

	return toMapping, nil
}

func ClientIP(c echo.Context) string {

	IPAddress := c.Request().Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = c.Request().Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = c.RealIP()
	}
	return IPAddress
}
