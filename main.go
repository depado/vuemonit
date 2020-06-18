package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/implem/auth.jwt"
	"github.com/Depado/vuemonit/implem/formatter.json"
	"github.com/Depado/vuemonit/implem/scheduler.gocron"
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

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Couldn't start")
	}
}
