package db

import (
	"github.com/SosisterRapStar/hotels/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const mysqlDriverName = "mysql"

// NewMySQL создаёт подключение к MySQL по конфигу.
func NewMySQL(cfg *config.Repository) (*sqlx.DB, error) {
	db, err := sqlx.Connect(mysqlDriverName, cfg.DSNMySQL())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxIdleTime(cfg.MaxIdleLifetime)
	db.SetConnMaxLifetime(cfg.MaxOpenLifetime)
	return db, nil
}
