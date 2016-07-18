package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	mw "github.com/kirikami/pomodoro_crew/middleware"
	a "github.com/kirikami/pomodoro_crew/utils/application"
)

func StartServer(application a.Application) {
	server := echo.New()

	server.Use(middleware.Recover())
	server.Use(mw.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin},
	}))

	authConfig := auth.ApiV1Handler{application.Configuration}

	server.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Status ok")
	})

	v1 := server.Group("/v1")

	aut := v1.Group("/auth")
	aut.GET("", authConfig.AutenteficationHandler)
	aut.GET("/callback", authConfig.ProviderCallback)

	server.Run(standard.New(fmt.Sprintf(":%d", application.Configuration.ListenAddress)))
}
