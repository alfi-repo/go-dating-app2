package storage

import (
	"database/sql"
	"go-dating-app/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.DB.MaxOpenPool)
	db.SetMaxIdleConns(config.DB.MaxIdlePool)
	db.SetConnMaxLifetime(time.Duration(config.DB.MaxIdleSecond) * time.Second)
	return db, nil
}
