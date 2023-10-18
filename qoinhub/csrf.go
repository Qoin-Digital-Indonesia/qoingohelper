/*
this middleware function is to validate csrf token is exist in redis
if exist continue to process the request
if not exist return error
*/

package qoingohelper

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func VerifCsrf(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var response ResponseApi

		header := &Headers{
			Csrf: c.Request().Header.Get("X-Xsrf-Token"),
		}
		// validate header request
		validate := header.ValiadateHeaderCsrf()
		if validate != nil {
			response.BadRequest("invalid validation", validate)
			return c.JSON(response.Code, response)
		}

		// get token from redis
		_, err := GetRedis("csrf-" + header.Csrf)
		if err != nil {
			// if error redis keys not found, return unathorized
			switch strings.Contains(err.Error(), "redis: nil") {
			case true:
				LoggerErrorHub(err)
				response.Unauthorized("Token invalid")
				return c.JSON(response.Code, response)
			case false:
				response.InternalServerError(err)
				return c.JSON(response.Code, response)
			}
		}

		return next(c)
	}
}
