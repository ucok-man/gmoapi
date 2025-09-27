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
	Port int         `env:"GMOAPI_PORT" envDefault:"4000"`
	Env  Environment `env:"GMOAPI_ENV" envDefault:"development"`
	DB   DatabaseConfig
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
	flag.IntVar(&cfg.Port, "port", cfg.Port, "API server port (1-65535)")
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
		printUsage()
		os.Exit(1)
	}
	return cfg
}

// printUsage prints usage information and available flags
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Configuration can be provided via environment variables or command-line flags.\n")
	fmt.Fprintf(os.Stderr, "Command-line flags override environment variables.\n\n")
	fmt.Fprintf(os.Stderr, "Environment Variables:\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_PORT                   - API server port (default: 4000)\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_ENV                    - Environment (default: development)\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_DB_DSN                 - PostgreSQL connection string (required)\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_DB_MAX_OPEN_CONN       - Max open database connections (default: 25)\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_DB_MAX_IDLE_CONN       - Max idle database connections (default: 25)\n")
	fmt.Fprintf(os.Stderr, "  GMOAPI_DB_MAX_IDLE_TIME       - Max database connection idle time (default: 15m)\n\n")
	fmt.Fprintf(os.Stderr, "Command-line Flags:\n")
	flag.PrintDefaults()
}
