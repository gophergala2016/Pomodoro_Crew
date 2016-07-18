package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"net/http"

	m "github.com/kirikami/go_exercise_api/database/models"
	u "github.com/kirikami/go_exercise_api/utils"
)

func (h ApiV1Handler) GetAllTasksHendler(c echo.Context) error {
	tasks := m.TaskList{}

	if err := h.Database.Find(&tasks.Tasks).Error; err != nil {
		u.SendError(http.StatusInternalServerError, c, err, DatabaseErrorMessage)
	}

	return c.JSON(http.StatusOK, tasks)
}
