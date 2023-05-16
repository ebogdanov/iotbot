package data

import (
	"database/sql"
	"fmt"
	"su27bot/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	Conn *sql.DB
}

func NewDb(cfg *config.Config) (*Db, error) {
	if cfg.Db == nil {
		return nil, fmt.Errorf("config file do not have section named %s", "db")
	}

	if cfg.Db.Driver == "" {
		return nil, fmt.Errorf("config file do not have section named %s", "db.driver")
	}

	if cfg.Db.Driver == "sqlite3" || cfg.Db.Driver == "mysql" || cfg.Db.Driver == "postgre" {
		db, err := sql.Open(cfg.Db.Driver, cfg.Db.Connection)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}

		return &Db{Conn: db}, err
	}

	return nil, fmt.Errorf("%s is not in supported engines list: (sqlite, mysql, postgre)", cfg.Db.Driver)
}
