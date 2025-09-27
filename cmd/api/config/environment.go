package config

import (
	"errors"
	"slices"
	"strings"
)

// Environment represents the application environment
type Environment string

// Valid environments
const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsValid checks if the environment is valid
func (e Environment) IsValid() bool {
	validEnvs := []Environment{EnvDevelopment, EnvStaging, EnvProduction}
	return slices.Contains(validEnvs, e)
}

// IsDevelopment returns true if the environment is development
func (e Environment) IsDevelopment() bool {
	return e == EnvDevelopment
}

// IsProduction returns true if the environment is production
func (e Environment) IsProduction() bool {
	return e == EnvProduction
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (e *Environment) UnmarshalText(value []byte) error {
	input := strings.TrimSpace(strings.ToLower(string(value)))
	env := Environment(input)

	if !env.IsValid() {
		return errors.New("invalid environment: must be one of (development|staging|production)")
	}

	*e = env
	return nil
}
