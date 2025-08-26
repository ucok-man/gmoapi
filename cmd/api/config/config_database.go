package config

import (
	"fmt"
	"time"
)

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Password    string        `env:"GMOAPI_DB_PASSWORD"`
	MaxOpenConn int           `env:"GMOAPI_DB_MAX_OPEN_CONN" envDefault:"25"`
	MaxIdleConn int           `env:"GMOAPI_DB_MAX_IDLE_CONN" envDefault:"25"`
	MaxIdleTime time.Duration `env:"GMOAPI_DB_MAX_IDLE_TIME" envDefault:"15m"`
}

// Validate validates the database configuration
func (d *DatabaseConfig) Validate() error {
	if d.Password == "" {
		return ErrMissingDBPassword
	}

	if d.MaxOpenConn < 1 {
		return fmt.Errorf("%w: max open connections must be at least 1", ErrInvalidDBConfig)
	}

	if d.MaxIdleConn < 0 {
		return fmt.Errorf("%w: max idle connections cannot be negative", ErrInvalidDBConfig)
	}

	if d.MaxIdleConn > d.MaxOpenConn {
		return fmt.Errorf("%w: max idle connections cannot exceed max open connections", ErrInvalidDBConfig)
	}

	if d.MaxIdleTime < 0 {
		return fmt.Errorf("%w: max idle time cannot be negative", ErrInvalidDBConfig)
	}

	return nil
}
