package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/implem/auth.jwt"
	"github.com/Depado/vuemonit/implem/formatter.json"
	"github.com/Depado/vuemonit/implem/scheduler.gocron"
	noopsched "github.com/Depado/vuemonit/implem/scheduler.noop"
	"github.com/Depado/vuemonit/implem/storage.storm"
	"github.com/Depado/vuemonit/infra"
	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/router"
)

// Build number and versions injected at compile time, set yours
var (
	Version = "unknown"
	Build   = "unknown"
)

// Main command that will be run when no other command is provided on the
// command-line
var rootCmd = &cobra.Command{
	Use: "vuemonit",
	Run: func(cmd *cobra.Command, args []string) { run() },
}

// Version command that will display the build number and version (if any)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Printf("Build: %s\nVersion: %s\n", Build, Version) },
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "manage users",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("hello") },
}

var addUser = &cobra.Command{
	Use:   "add [email] [password]",
	Short: "add a user",
	Args:  cobra.MinimumNArgs(2),
	Run: func(c *cobra.Command, args []string) {
		type Params struct {
			Email    string
			Password string
		}
		params := func() *Params {
			return &Params{Email: args[0], Password: args[1]}
		}
		create := func(p *Params, l *zerolog.Logger, i interactor.LogicHandler) {
			if err := i.Register(p.Email, p.Password); err != nil {
				l.Err(err).Msg("unable to create new user")
			} else {
				l.Info().Str("email", p.Email).Msg("created user")
			}
		}
		app := fx.New(
			fx.NopLogger,
			fx.Provide(
				cmd.NewConf,
				cmd.NewLogger,
				storage.NewStormStorage,
				noopsched.NewNoopScheduler,
				auth.NewJWTAuthProvider,
				formatter.NewJSONFormatter,
				interactor.NewInteractor,
				params,
			),
			fx.Invoke(create),
		)
		app.Start(context.Background()) // nolint: errcheck
		app.Stop(context.Background())  // nolint: errcheck
	},
}

func run() {
	fx.New(
		fx.Provide(
			cmd.NewConf,
			cmd.NewLogger,
			infra.NewCORS,
			infra.NewServer,
			storage.NewStormStorage,
			scheduler.NewGocronScheduler,
			auth.NewJWTAuthProvider,
			formatter.NewJSONFormatter,
			interactor.NewInteractor,
			router.NewRouter,
		),
		fx.Invoke(router.Run),
	).Run()
}

func main() {
	// Initialize Cobra and Viper
	// cobra.OnInitialize(cmd.Initialize)
	cmd.AddAllFlags(rootCmd)
	rootCmd.AddCommand(versionCmd)
	userCmd.AddCommand(addUser)
	rootCmd.AddCommand(userCmd)

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Couldn't start")
	}
}
