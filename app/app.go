package app

import (
	"GoWeb/config"
	"database/sql"
	"embed"
)

type App struct {
	Config         config.Configuration
	Db             *sql.DB
	Res            *embed.FS
	ScheduledTasks Scheduled
}
