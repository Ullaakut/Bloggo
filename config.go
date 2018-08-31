package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Config represents the Bloggo configuration
type Config struct {
	LogLevel      string `json:"log_level"`
	ServerAddress string `json:"server_address"`
	ServerPort    uint   `json:"server_port"`

	MySQLURL           string        `json:"mysql_url"`
	MySQLRetryInterval time.Duration `json:"mysql_retry_interval"`
	MySQLRetryDuration time.Duration `json:"mysql_retry_duration"`

	jwtSecret string
}

// Set default values for configuration parameters
func init() {
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("server_address", "0.0.0.0")
	viper.SetDefault("server_port", 4242)
	viper.SetDefault("mysql_url", "root:root@tcp(db:3306)/bloggo?charset=utf8&parseTime=True&loc=Local")
	viper.SetDefault("mysql_retry_interval", "2s")
	viper.SetDefault("mysql_retry_duration", "1m")
}

// GetConfig sets the default values for the configuration and gets it from the environment/command line
func GetConfig() Config {
	var config Config

	// Override default with environment variables
	viper.SetEnvPrefix("BLOGGO")
	viper.AutomaticEnv()
	viper.Unmarshal(&config)

	config.jwtSecret = viper.GetString("jwt_secret")

	config.LogLevel = viper.GetString("log_level")
	config.ServerAddress = viper.GetString("server_address")
	config.ServerPort = uint(viper.GetInt("server_port"))
	config.MySQLURL = viper.GetString("mysql_url")

	config.MySQLRetryInterval = viper.GetDuration("mysql_retry_interval")
	config.MySQLRetryDuration = viper.GetDuration("mysql_retry_duration")

	return config
}

// Print prints the current configuration
func (c Config) Print(log *zerolog.Logger) {
	log.Debug().
		Str("LogLevel", c.LogLevel).
		Str("ServerAddress", c.ServerAddress).
		Uint("ServerPort", c.ServerPort).
		Str("MySQLURL", c.MySQLURL).
		Dur("MySQLRetryInterval", c.MySQLRetryInterval).
		Dur("MySQLRetryDuration", c.MySQLRetryDuration).
		Msg("Configuration")
}
