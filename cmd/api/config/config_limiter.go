package config

import "errors"

type LimiterConfig struct {
	Rps     float64 `env:"GMOAPI_LIMITER_RPS" envDefault:"2"`
	Burst   int     `env:"GMOAPI_LIMITER_BURST" envDefault:"4"`
	Enabled bool    `env:"GMOAPI_LIMITER_ENABLED" envDefault:"true"`
}

func (c *LimiterConfig) Validate() error {
	if c.Enabled {
		if c.Rps < 0 {
			return errors.New("limiter rps must be positive")
		}

		if c.Burst < 0 {
			return errors.New("limiter burst must be positive")
		}
	}

	return nil
}
