package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/markbates/goth/gothic"

	_ "github.com/kirikami/go_exercise_api/config"
)

type TokenMessage struct {
	Token string `json:"token"`
}

func (h *ApiV1Handler) ProviderCallback(c echo.Context) error {

	res := c.Response().(*standard.Response).ResponseWriter
	req := c.Request().(*standard.Request).Request
	_, err := gothic.CompleteUserAuth(res, req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway)
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(h.Configuration.SigningKey)
	sentToken := TokenMessage{t}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sentToken)
}
