package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	v "gopkg.in/go-playground/validator.v9"
)

// Config represents the Bloggo configuration
type Config struct {
	LogLevel      string `json:"log_level" validate:"required,eq=DEBUG|eq=INFO|eq=WARNING|eq=ERROR|eq=FATAL"`
	ServerAddress string `json:"server_address" validate:"required"`
	ServerPort    uint   `json:"server_port" validate:"required,min=1,max=65535"`

	MySQLURL           string        `json:"mysql_url"`
	MySQLRetryInterval time.Duration `json:"mysql_retry_interval"`
	MySQLRetryDuration time.Duration `json:"mysql_retry_duration"`

	BcryptRuns int `json:"bcrypt_runs" validate:"min=4,max=31"`

	JWTSecret string `validate:"required,min=1"`
}

// Set default values for configuration parameters
func init() {
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("server_address", "0.0.0.0")
	viper.SetDefault("server_port", 4242)
	viper.SetDefault("mysql_url", "root:root@tcp(db:3306)/bloggo?charset=utf8&parseTime=True&loc=Local")
	viper.SetDefault("mysql_retry_interval", "2s")
	viper.SetDefault("mysql_retry_duration", "1m")
	viper.SetDefault("bcrypt_runs", 11)
}

// GetConfig sets the default values for the configuration and gets it from the environment/command line
func GetConfig() (Config, error) {
	var config Config

	// Override default with environment variables
	viper.SetEnvPrefix("BLOGGO")
	viper.AutomaticEnv()
	viper.Unmarshal(&config)

	config.JWTSecret = viper.GetString("jwt_secret")

	config.LogLevel = viper.GetString("log_level")
	config.ServerAddress = viper.GetString("server_address")
	config.ServerPort = uint(viper.GetInt("server_port"))
	config.MySQLURL = viper.GetString("mysql_url")

	config.MySQLRetryInterval = viper.GetDuration("mysql_retry_interval")
	config.MySQLRetryDuration = viper.GetDuration("mysql_retry_duration")

	config.BcryptRuns = viper.GetInt("bcrypt_runs")

	validate := v.New()
	err := validate.Struct(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Print prints the current configuration
func (c Config) Print(log *zerolog.Logger) {
	log.Debug().
		Str("log_level", c.LogLevel).
		Str("server_address", c.ServerAddress).
		Uint("server_port", c.ServerPort).
		Str("mysql_url", c.MySQLURL).
		Dur("mysql_retry_interval", c.MySQLRetryInterval).
		Dur("mysql_retry_duration", c.MySQLRetryDuration).
		Int("bcrypt_runs", c.BcryptRuns).
		Msg("configuration")
}
