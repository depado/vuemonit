package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConf struct {
	Level  string `mapstructure:"level"`
	Preset string `mapstructure:"preset"`
}

type ServerConf struct {
	Host string   `mapstructure:"host"`
	Port int      `mapstructure:"port"`
	Mode string   `mapstructure:"mode"`
	Log  bool     `mapstructure:"log"`
	CORS CorsConf `mapstructure:"cors"`
}

type CorsConf struct {
	Enable  bool     `mapstructure:"enable"`
	Methods []string `mapstructure:"methods"`
	Expose  []string `mapstructure:"expose"`
	Headers []string `mapstructure:"headers"`
	Origins []string `mapstructure:"origins"`
	All     bool     `mapstructure:"all"`
}

type PrometheusConf struct {
	Prefix   string `mapstructure:"prefix"`
	Disabled bool   `mapstructure:"disabled"`
}

type Conf struct {
	Log        LogConf        `mapstructure:"log"`
	Server     ServerConf     `mapstructure:"server"`
	Prometheus PrometheusConf `mapstructure:"prometheus"`
}

// NewLogger will return a new logger
func NewLogger(c *Conf) (*zap.Logger, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(c.Log.Level)); err != nil {
		return nil, fmt.Errorf("parse level: %w", err)
	}
	switch c.Log.Preset {
	case "development":
		config := zap.NewDevelopmentConfig()
		config.Level = lvl
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := config.Build()
		if err != nil {
			return logger, fmt.Errorf("unable to create dev config: %w", err)
		}
		return logger, err
	case "production":
		config := zap.NewProductionConfig()
		config.Level = lvl
		return config.Build()
	default:
		return nil, fmt.Errorf("unrecognized logger preset: %s", c.Log.Preset)
	}
}

// NewConf will parse and return the configuration
func NewConf() (*Conf, error) {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("vuemonit")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Configuration file
	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}

	viper.ReadInConfig() // nolint: errcheck
	conf := &Conf{}
	if err := viper.Unmarshal(conf); err != nil {
		return conf, fmt.Errorf("unable to unmarshal conf: %w", err)
	}

	return conf, nil
}
