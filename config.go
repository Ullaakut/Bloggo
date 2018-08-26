package main

import (
	"time"

	"github.com/rs/zerolog"
)

// Config represents the Bloggo configuration
type Config struct {
	LogLevel                string
	GracefulShutdownTimeout time.Duration
	ServerAddress           string
	ServerPort              uint

	TrustedSource string
	JWTSecret     string

	MySQLURL           string
	MySQLRetryInterval time.Duration
	MySQLRetryDuration time.Duration
}

// DefaultConfig generates a configuration structure with the default values
func DefaultConfig() Config {
	return Config{
		LogLevel:                "DEBUG",
		GracefulShutdownTimeout: 10 * time.Second,
		ServerAddress:           "0.0.0.0",
		ServerPort:              4242,

		TrustedSource: "src.bloggo.com",

		MySQLURL:           "root:root@tcp(db:3306)/bloggo?charset=utf8&parseTime=True&loc=Local",
		MySQLRetryInterval: 2 * time.Second,
		MySQLRetryDuration: 1 * time.Minute,
	}
}

// Print prints the current configuration
func (c Config) Print(log *zerolog.Logger) {
	log.Debug().
		Str("LogLevel", c.LogLevel).
		Str("ServerAddress", c.ServerAddress).
		Uint("ServerPort", c.ServerPort).
		Dur("GraciousShutdownTimeout", c.GracefulShutdownTimeout).
		Str("MySQLURL", c.MySQLURL).
		Dur("MySQLRetryInterval", c.MySQLRetryInterval).
		Dur("MySQLRetryDuration", c.MySQLRetryDuration).
		Str("TrustedSource", c.TrustedSource).
		Str("JWTSecret", "**************").
		Msg("Configuration")
}
