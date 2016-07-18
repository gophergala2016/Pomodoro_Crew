package application

import (
	"flag"

	"github.com/google/cayley"
	"github.com/google/cayley/graph"

	"github.com/kirikami/pomodoro_crew/config"
	"github.com/kirikami/pomodoro_crew/database"
)

type App interface {
	InitConfiguration()
	InitDatabase()
}

type Application struct {
	Configuration *config.Configuration
	Database      *database.Storage
}

func (a *Application) InitConfiguration() {
	configfile := flag.String("config", "config.json", "Config for connection to database")
	flag.Parse()
	a.Configuration = config.MustNewConfig(*configfile)
}

func (a *Application) InitDatabase() {
	a.Database = database.MustNewDatabase(a.Configuration.DatabaseConfig)
}
