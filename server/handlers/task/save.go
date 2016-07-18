package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"net/http"

	m "github.com/kirikami/go_exercise_api/database/models"
	u "github.com/kirikami/go_exercise_api/utils"
)

func (h ApiV1Handler) SaveTaskHandler(c echo.Context) error {
	task := m.Task{}

	if err := c.Bind(&task); err != nil {
		u.SendError(StatusUnprocessableEntity, c, err, IncorrectDataErrorMessage)
	}

	if err := h.Database.Save(&task).Error; err != nil {
		u.SendError(http.StatusInternalServerError, c, err, DatabaseErrorMessage)

	}

	return c.JSON(http.StatusCreated, task)
}
