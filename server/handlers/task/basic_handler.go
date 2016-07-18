package task

import (
	"github.com/jinzhu/gorm"
	"github.com/kirikami/go_exercise_api/config"
)

type ApiV1Handler struct {
	Configuration *config.Configuration
	Database      *gorm.DB
}

var (
	IdErrorMessage            string = "Id should be numeric and exist in database"
	DatabaseErrorMessage      string = "Database internal error"
	IncorrectDataErrorMessage string = "Incorrect data or format of data sent"
	StatusUnprocessableEntity int    = 422
)
