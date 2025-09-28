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

func (c *DatabaseConfig) Validate() error {
	if c.DSN == "" {
		return errors.New("database connection string is required")
	}

	if c.MaxOpenConn < 1 {
		return errors.New("max open connections must be at least 1")
	}

	if c.MaxIdleConn < 0 {
		return errors.New("max idle connections cannot be negative")
	}

	if c.MaxIdleConn > c.MaxOpenConn {
		return errors.New("max idle connections cannot exceed max open connections")
	}

	if c.MaxIdleTime < 0 {
		return errors.New("max idle time cannot be negative")
	}

	return nil
}
