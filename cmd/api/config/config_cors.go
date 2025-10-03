package config

type CorsConfig struct {
	TrustedOrigins []string `env:"GMOAPI_CORS_TRUSTED_ORIGINS"`
}

func (c *CorsConfig) Validate() error {
	if len(c.TrustedOrigins) > 0 {

	}

	return nil
}
