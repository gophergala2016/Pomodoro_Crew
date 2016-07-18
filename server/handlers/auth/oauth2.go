package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	_ "github.com/kirikami/go_exercise_api/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
)

func (h *ApiV1Handler) AutenteficationHandler(c echo.Context) error {
	fb := h.Configuration.FacebookConfig
	goth.UseProviders(
		facebook.New(fb.ID, fb.Key, fb.CallbackAddress),
	)
	gothic.BeginAuthHandler(c.Response().(*standard.Response).ResponseWriter, c.Request().(*standard.Request).Request)

	return c.NoContent(http.StatusNoContent)
}
