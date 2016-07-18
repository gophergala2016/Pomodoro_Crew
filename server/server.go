package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	mw "github.com/kirikami/go_exercise_api/middleware"
	"github.com/kirikami/go_exercise_api/server/handlers/auth"
	t "github.com/kirikami/go_exercise_api/server/handlers/task"
	a "github.com/kirikami/go_exercise_api/utils/application"
)

func StartServer(application a.Application) {
	server := echo.New()

	server.Use(middleware.Recover())
	server.Use(mw.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin},
	}))

	taskConfig := t.ApiV1Handler{application.Configuration, application.Database}
	authConfig := auth.ApiV1Handler{application.Configuration}

	server.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Status ok")
	})

	v1 := server.Group("/v1")

	v1.GET("/tasks", taskConfig.GetAllTasksHendler, mw.JWT(taskConfig.Configuration.SigningKey))

	aut := v1.Group("/auth")
	aut.GET("", authConfig.AutenteficationHandler)
	aut.GET("/callback", authConfig.ProviderCallback)

	task := v1.Group("/task")
	task.Use(mw.JWT(taskConfig.Configuration.SigningKey))

	task.POST("", taskConfig.SaveTaskHandler)

	task.PUT("/:id", taskConfig.UpdateTaskHandler)
	task.DELETE("/:id", taskConfig.DeleteTaskHandler)
	task.GET("/:id", taskConfig.GetTaskHandler)

	server.Run(standard.New(fmt.Sprintf(":%d", application.Configuration.ListenAddress)))
}
