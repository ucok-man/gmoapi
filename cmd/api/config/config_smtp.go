package config

import (
	"fmt"
	"net/mail"
)

type SMTPConfig struct {
	Host     string `env:"GMOAPI_SMTP_HOST"`
	Port     int    `env:"GMOAPI_SMTP_PORT"`
	Username string `env:"GMOAPI_SMTP_USERNAME"`
	Password string `env:"GMOAPI_SMTP_PASSWORD"`
	Sender   string `env:"GMOAPI_SMTP_SENDER" envDefault:"Gmoapi <no-reply@gmoapi.ucokman.web.id>"`
}

func (c *SMTPConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("SMTP host must not be empty")
	}

	if c.Port <= 0 {
		return fmt.Errorf("SMTP port must be a positive number")
	}

	if c.Username == "" {
		return fmt.Errorf("SMTP username must not be empty")
	}

	if c.Password == "" {
		// In a real application, you might skip this check if the server
		// doesn't require a password (e.g., local mail server or specific setup).
		// However, for typical external SMTP, it's required.
		return fmt.Errorf("SMTP password must not be empty")
	}

	if _, err := mail.ParseAddress(c.Sender); err != nil {
		return fmt.Errorf("SMTP sender email is not a valid email address: %w", err)
	}

	return nil
}
