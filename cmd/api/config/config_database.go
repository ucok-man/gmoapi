package config

import (
	"errors"
	"time"
)

type DatabaseConfig struct {
	DSN         string        `env:"GMOAPI_DB_DSN"`
	MaxOpenConn int           `env:"GMOAPI_DB_MAX_OPEN_CONN" envDefault:"25"`
	MaxIdleConn int           `env:"GMOAPI_DB_MAX_IDLE_CONN" envDefault:"25"`
	MaxIdleTime time.Duration `env:"GMOAPI_DB_MAX_IDLE_TIME" envDefault:"15m"`
}

func (d *DatabaseConfig) Validate() error {
	if d.DSN == "" {
		return errors.New("database connection string is required")
	}

	if d.MaxOpenConn < 1 {
		return errors.New("max open connections must be at least 1")
	}

	if d.MaxIdleConn < 0 {
		return errors.New("max idle connections cannot be negative")
	}

	if d.MaxIdleConn > d.MaxOpenConn {
		return errors.New("max idle connections cannot exceed max open connections")
	}

	if d.MaxIdleTime < 0 {
		return errors.New("max idle time cannot be negative")
	}

	return nil
}
