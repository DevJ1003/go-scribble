package dbrepo

import (
	"database/sql"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/repository"
)

type mysqlDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewMysqlRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &mysqlDBRepo{
		App: a,
		DB:  conn,
	}
}
