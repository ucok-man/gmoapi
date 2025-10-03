package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

// Config represents the application configuration
type Config struct {
	Port    int         `env:"GMOAPI_PORT" envDefault:"4000"`
	Env     Environment `env:"GMOAPI_ENV" envDefault:"development"`
	DB      DatabaseConfig
	Limiter LimiterConfig
	SMTP    SMTPConfig
	Cors    CorsConfig
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	// Validate port
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	// Validate environment
	if !c.Env.IsValid() {
		return errors.New("invalid environment: must be one of (development|staging|production)")
	}

	// Validate database configuration
	if err := c.DB.Validate(); err != nil {
		return err
	}

	// Validate limiter configuration
	if err := c.Limiter.Validate(); err != nil {
		return err
	}

	// Validate cors configuration
	if err := c.Cors.Validate(); err != nil {
		return err
	}

	return nil
}

// NewConfig creates and validates a new configuration instance
func NewConfig() (Config, error) {
	cfg := Config{}

	// Parse environment variables
	if err := env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Define command-line flags
	flag.IntVar(&cfg.Port, "port", cfg.Port, "API server port")
	flag.Func("env", "Environment (development|staging|production)", func(s string) error {
		env := Environment(strings.TrimSpace(strings.ToLower(s)))
		if !env.IsValid() {
			return errors.New("invalid environment: must be one of (development|staging|production)")
		}
		cfg.Env = env
		return nil
	})

	flag.StringVar(&cfg.DB.DSN, "db-dsn", cfg.DB.DSN, "PostgreSQL connection string")
	flag.IntVar(&cfg.DB.MaxOpenConn, "db-max-open-conn", cfg.DB.MaxOpenConn, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConn, "db-max-idle-conn", cfg.DB.MaxIdleConn, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", cfg.DB.MaxIdleTime, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", cfg.Limiter.Rps, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", cfg.Limiter.Burst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", cfg.Limiter.Enabled, "Enable rate limiter")

	flag.StringVar(&cfg.SMTP.Host, "smtp-host", cfg.SMTP.Host, "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", cfg.SMTP.Port, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", cfg.SMTP.Username, "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", cfg.SMTP.Password, "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", cfg.SMTP.Sender, "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (comma separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Split(val, ",")
		return nil
	})

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Configuration can be provided via environment variables or command-line flags.\n")
		fmt.Fprintf(os.Stderr, "Command-line flags override environment variables.\n\n")
		fmt.Fprintf(os.Stderr, "Environment Variables:\n")
		fmt.Fprintf(os.Stderr, "  - To define environment variable use prefix GMOAPI.\n")
		fmt.Fprintf(os.Stderr, "  - The expected form is <PREFIX_FLAG_NAME> all in uppercase.\n")
		fmt.Fprintf(os.Stderr, "  - ex: GMOAPI_PORT, GMOAPI_DB_MAX_OPEN_CONN.\n\n")

		fmt.Fprintf(os.Stderr, "Command-line Flags:\n")
		flag.PrintDefaults()
	}

	// Parse command-line flags
	flag.Parse()

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// MustNewConfig creates a new configuration instance and exits on error
func MustNewConfig() Config {
	cfg, err := NewConfig()
	if err != nil {
		log.Printf("Configuration error: %v", err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	return cfg
}
