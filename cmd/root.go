package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddLoggerFlags adds support to configure the level of the logger.
func AddLoggerFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.type", "console", `one of "console" or "json"`)
	c.PersistentFlags().Bool("log.caller", true, "display the file and line where the call was made")
}

// AddServerFlags adds support to configure the server.
func AddServerFlags(c *cobra.Command) {
	// Server related flags
	c.PersistentFlags().String("server.host", "127.0.0.1", "host on which the server should listen")
	c.PersistentFlags().Int("server.port", 8080, "port on which the server should listen")
	c.PersistentFlags().String("server.mode", "release", "server mode can be either 'debug', 'test' or 'release'")
	c.PersistentFlags().Bool("server.log", true, "log requests or not")

	// CORS related flags
	c.PersistentFlags().Bool("server.cors.enable", false, "enable CORS")
	c.PersistentFlags().StringSlice("server.cors.methods", []string{"GET", "PUT", "POST", "DELETE", "OPTION", "PATCH"}, "array of allowed method when cors is enabled")
	c.PersistentFlags().StringSlice("server.cors.headers", []string{"Origin", "Authorization", "Content-Type"}, "array of allowed headers")
	c.PersistentFlags().StringSlice("server.cors.expose", []string{}, "array of exposed headers")
	c.PersistentFlags().StringSlice("server.cors.origins", []string{}, "array of allowed origins (overwritten if all is active)")
	c.PersistentFlags().Bool("server.cors.all", false, "defines that all origins are allowed")

	// JWT related flags
	c.PersistentFlags().String("server.jwt.secret", "", "secret to generate JWT token")
	c.PersistentFlags().String("server.cookie.domain", "", "domain on which the cookie might be used")
	c.PersistentFlags().Bool("server.cookie.https", true, "always send the cookie to https endpoints")

}

func AddFrontFlags(c *cobra.Command) {
	c.PersistentFlags().Bool("front.serve", true, "let the server serve the frontend")
	c.PersistentFlags().String("front.path", "front/dist/spa", "path to the frontend build")
}

func AddDatabaseFlags(c *cobra.Command) {
	c.PersistentFlags().String("database.path", "monit.db", "path to the database file to use")
}

// AddPrometheusFlags adds flags to support prometheus instrumentation.
func AddPrometheusFlags(c *cobra.Command) {
	c.PersistentFlags().String("prometheus.prefix", "vuemonit", "prefix for the prometheus label")
	c.PersistentFlags().Bool("prometheus.disabled", false, "enable or disable prometheus instrumentation")
}

// AddConfigurationFlag adds support to provide a configuration file on the
// command line.
func AddConfigurationFlag(c *cobra.Command) {
	c.PersistentFlags().String("conf", "", "configuration file to use")
}

// AddAllFlags will add all the flags provided in this package to the provided
// command and will bind those flags with viper.
func AddAllFlags(c *cobra.Command) {
	AddConfigurationFlag(c)
	AddLoggerFlags(c)
	AddPrometheusFlags(c)
	AddServerFlags(c)
	AddFrontFlags(c)
	AddDatabaseFlags(c)

	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddAllFlags").Fatal("Couldn't bind flags")
	}
}
