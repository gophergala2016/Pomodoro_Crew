package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"net/http"

	m "github.com/kirikami/go_exercise_api/database/models"
	u "github.com/kirikami/go_exercise_api/utils"
)

func (h ApiV1Handler) UpdateTaskHandler(c echo.Context) error {
	id, err := u.ParseIdInt64FromString(c.Param("id"))

	if err != nil {
		u.SendError(http.StatusBadRequest, c, err, IdErrorMessage)

	}

	task := m.Task{}

	if err := h.Database.First(&task, id).Error; err != nil {
		u.SendError(http.StatusInternalServerError, c, err, DatabaseErrorMessage)

	}

	if err := c.Bind(&task); err != nil {
		u.SendError(StatusUnprocessableEntity, c, err, IncorrectDataErrorMessage)
	}

	task.SetIsCompleted()

	if err := h.Database.Save(&task).Error; err != nil {
		u.SendError(http.StatusInternalServerError, c, err, DatabaseErrorMessage)
	}

	return c.JSON(http.StatusOK, task)
}
