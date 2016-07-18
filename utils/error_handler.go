package utils

import (
	"github.com/labstack/echo"
)

type SendErrorMessage struct {
	Message string `json:"message"`
}

func SendError(code int, c echo.Context, err error, msg string) {
	message := SendErrorMessage{msg}
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if !c.Response().Committed() {
		c.JSON(code, message)
	}
	return
}
