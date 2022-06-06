package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mariaiu/yandex-eda-app/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName, cfg.DB.Password, cfg.DB.SSLMode))
	if err != nil {
			return nil, err
	}
	err = db.Ping(); if err != nil {
		return nil, err
	}

	return db, err
}