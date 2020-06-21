package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type LogConf struct {
	Level  string `mapstructure:"level"`
	Type   string `mapstructure:"type"`
	Caller bool   `mapstructure:"caller"`
}

type JWTConf struct {
	Secret string `mapstructure:"secret"`
}

type CookieConf struct {
	Domain string `mapstructure:"domain"`
	HTTPS  bool   `mapstructure:"https"`
}

type ServerConf struct {
	Host   string     `mapstructure:"host"`
	Port   int        `mapstructure:"port"`
	Mode   string     `mapstructure:"mode"`
	Log    bool       `mapstructure:"log"`
	CORS   CorsConf   `mapstructure:"cors"`
	JWT    JWTConf    `mapstructure:"jwt"`
	Cookie CookieConf `mapstructure:"cookie"`
}

type FrontConf struct {
	Serve bool   `mapstructure:"serve"`
	Path  string `mapstructure:"path"`
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

type DatabaseConf struct {
	Path string `mapstructure:"path"`
}

type Conf struct {
	Log        LogConf        `mapstructure:"log"`
	Server     ServerConf     `mapstructure:"server"`
	Front      FrontConf      `mapstructure:"front"`
	Prometheus PrometheusConf `mapstructure:"prometheus"`
	Database   DatabaseConf   `mapstructure:"database"`
	Register   bool           `mapstructure:"allow_register"`
}

// NewLogger will return a new logger
func NewLogger(c *Conf) *zerolog.Logger {
	// Level parsing
	warns := []string{}
	lvl, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		warns = append(warns, fmt.Sprintf("unrecognized log level '%s', fallback to 'info'", c.Log.Level))
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(lvl)
	}

	// Type parsing
	switch c.Log.Type {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "json":
		break
	default:
		warns = append(warns, fmt.Sprintf("unrecognized log type '%s', fallback to 'json'", c.Log.Type))
	}

	// Caller
	if c.Log.Caller {
		log.Logger = log.With().Caller().Logger()
	}

	// Log messages with the newly created logger
	for _, w := range warns {
		log.Warn().Msg(w)
	}

	return &log.Logger
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
