package qoingohelper

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int         `json:"status_code"`
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type V2Response struct {
	Code    int         `json:"StatusCode"`
	Status  string      `json:"Status"`
	Message interface{} `json:"Message"`
	Data    interface{} `json:"Data"`
}

type JSONResponse struct {
	Code    int         `json:"status_code"`
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type V2JSONResponse struct {
	Code    int         `json:"StatusCode"`
	Status  string      `json:"Status"`
	Message interface{} `json:"Message"`
	Data    interface{} `json:"Data"`
}

// V1 Version 1 snake_case
func SuccessContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status_code": 200,
		"status":      "success",
		"message":     message,
		"data":        data,
	})
}

func ErrorContext(message interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status_code": 500,
		"status":      "failed",
		"message":     message,
		"data":        nil,
	})
}

func ValidationContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status_code": 400,
		"status":      "validation",
		"message":     message,
		"data":        data,
	})
}

func TimeoutContext(message interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status_code": 504,
		"status":      "timeout",
		"message":     message,
		"data":        nil,
	})
}

func NotFoundContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"status_code": 404,
		"status":      "not found",
		"message":     message,
		"data":        nil,
	})
}

func ResponseContext(code int, message interface{}, data interface{}, c echo.Context) error {
	if code == 200 { // Success
		return SuccessContext(message, data, c)
	} else if code == 400 { // Bad Request
		return ValidationContext(message, data, c)
	} else if code == 404 { // Notfound
		return NotFoundContext(message, data, c)
	} else if code == 504 { // Timeout
		return TimeoutContext(message, c)
	}
	return ErrorContext(message, c) // Internal Server Error
}

func ValidationResp(message interface{}, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}

func Success(message string, data interface{}, c echo.Context) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}

func (response *Response) Success(code int, message string, data interface{}) {
	response.Code = code
	response.Status = "success"
	response.Message = message
	response.Data = data
}

func (response *Response) Error(code int, message string, data interface{}) {
	response.Code = code
	response.Status = "error"
	response.Message = message
	response.Data = data
}

func (response *JSONResponse) Success(message interface{}, data interface{}) {
	response.Status = "success"
	response.Message = message
	response.Data = data
	response.Code = http.StatusOK
}

func (response *JSONResponse) Error(message interface{}, data interface{}) {
	response.Status = "error"
	response.Message = message
	response.Data = data
	response.Code = http.StatusInternalServerError
}

func (response *JSONResponse) Validation(message interface{}, data interface{}) {
	response.Status = "validation"
	response.Message = message
	response.Data = data
	response.Code = http.StatusBadRequest
}

func RPCJSONResponse(status string, message interface{}, data interface{}) string {
	var responseStruct = new(JSONResponse)

	if status == "success" {
		responseStruct.Success(message, data)
	} else if status == "validation" {
		responseStruct.Validation(message, data)
	} else {
		if data == nil {
			responseStruct.Error(message, nil)
		} else if fmt.Sprintf("%v", reflect.TypeOf(data).Kind()) == "ptr" {
			responseStruct.Error(message, fmt.Sprintf("%v", data))
		} else {
			responseStruct.Error(message, data)
		}
	}

	return JsonEncode(responseStruct)
}

// V2 Version 2 CamlCase
func V2SuccessContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"StatusCode": 200,
		"Status":     "success",
		"Message":    message,
		"Data":       data,
	})
}

func V2ErrorContext(message interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"StatusCode": 500,
		"Status":     "failed",
		"Message":    message,
		"Data":       nil,
	})
}

func V2ValidationContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"StatusCode": 400,
		"Status":     "validation",
		"Message":    message,
		"Data":       data,
	})
}

func V2TimeoutContext(message interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"StatusCode": 504,
		"Status":     "timeout",
		"Message":    message,
		"Data":       nil,
	})
}

func V2NotFoundContext(message interface{}, data interface{}, c echo.Context) (err error) {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"StatusCode": 404,
		"Status":     "not found",
		"Message":    message,
		"Data":       nil,
	})
}

func V2ResponseContext(code int, message interface{}, data interface{}, c echo.Context) error {
	if code == 200 { // Success
		return SuccessContext(message, data, c)
	} else if code == 400 { // Bad Request
		return ValidationContext(message, data, c)
	} else if code == 404 { // Notfound
		return NotFoundContext(message, data, c)
	} else if code == 504 { // Timeout
		return TimeoutContext(message, c)
	}
	return ErrorContext(message, c) // Internal Server Error
}

func V2ValidationResp(message interface{}, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Status":  "success",
		"Message": message,
		"Data":    data,
	}
}

func V2Success(message string, data interface{}, c echo.Context) map[string]interface{} {
	return map[string]interface{}{
		"Status":  "success",
		"Message": message,
		"Data":    data,
	}
}

func (response *V2Response) V2Success(code int, message string, data interface{}) {
	response.Code = code
	response.Status = "success"
	response.Message = message
	response.Data = data
}

func (response *V2Response) V2Error(code int, message string, data interface{}) {
	response.Code = code
	response.Status = "error"
	response.Message = message
	response.Data = data
}

func (response *V2JSONResponse) V2Success(message interface{}, data interface{}) {
	response.Status = "success"
	response.Message = message
	response.Data = data
	response.Code = http.StatusOK
}

func (response *V2JSONResponse) V2Error(message interface{}, data interface{}) {
	response.Status = "error"
	response.Message = message
	response.Data = data
	response.Code = http.StatusInternalServerError
}

func (response *V2JSONResponse) V2Validation(message interface{}, data interface{}) {
	response.Status = "validation"
	response.Message = message
	response.Data = data
	response.Code = http.StatusBadRequest
}

func V2RPCJSONResponse(status string, message interface{}, data interface{}) string {
	var responseStruct = new(V2JSONResponse)

	if status == "success" {
		responseStruct.V2Success(message, data)
	} else if status == "validation" {
		responseStruct.V2Validation(message, data)
	} else {
		if data == nil {
			responseStruct.V2Error(message, nil)
		} else if fmt.Sprintf("%v", reflect.TypeOf(data).Kind()) == "ptr" {
			responseStruct.V2Error(message, fmt.Sprintf("%v", data))
		} else {
			responseStruct.V2Error(message, data)
		}
	}

	return JsonEncode(responseStruct)
}
