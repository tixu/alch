package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"

	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
)

// MainConfigPath Configuration path
const MainConfigPath = "conf/config.yaml"

var k = koanf.New(".")
var validate = validator.New()

type Config struct {
	Log    *LogConfig    `koanf:"log"`
	Cachet *CachetConfig `koanf:"cachet" validate:"required"`
	Server *ServerConfig `koanf:"server"`
}

// CachetConfig CachetHQ Configuration
type CachetConfig struct {
	URL    string `koanf:"url" validate:"required,uri"`
	APIKey string `koanf:"apiKey" validate:"required"`
}

type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}

type ServerConfig struct {
	ListenAddr string `koanf:"listenAddr"`
	Port       int    `koanf:"port" validate:"required"`
}

func Load() (*Config, error) {

	// Try to load main configuration file
	err := k.Load(file.Provider(MainConfigPath), yaml.Parser())
	if err != nil {
		return nil, err
	}

	// Prepare configuration object
	var out Config
	// Quick unmarshal.
	err = k.Unmarshal("", &out)
	if err != nil {
		return nil, err
	}

	// Configuration validation
	err = validate.Struct(out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// ConfigureLogger Configure logger instance
func ConfigureLogger(logger *logrus.Logger, logConfig *LogConfig) error {
	// Manage log format
	if logConfig.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// Manage log level
	lvl, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		return err
	}
	logger.SetLevel(lvl)

	return nil
}
