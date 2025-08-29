package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

func InitDB(cfg DSN) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
