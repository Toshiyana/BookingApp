package dbrepo

import (
	"database/sql"

	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// If you want to use other kinds of databases, you can create functions.
// ex: func NewMysqlRepo()
func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
