package main

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

//go:embed db/*.sql
var dbDir embed.FS

var (
	dsnURI string

	logLevelArg string
	logLevel    = &slog.LevelVar{}
	logger      = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}))

	cmdRoot = &cobra.Command{
		Use:   "mcp",
		Short: "MCP is a cli tool for managing all of your model context needs via the Model Context Protocol.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	cmdRoot.PersistentFlags().StringVar(&logLevelArg, "log-level", "info", "log level among \"debug\", \"info\" or \"error\"")

	cmdRoot.AddCommand(cmdServe)
}

func initConfig() {
	switch logLevelArg {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "error":
		logLevel.Set(slog.LevelError)
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	cfgDir := path.Join(home, ".mcp")

	cobra.CheckErr(os.MkdirAll(cfgDir, 0750))

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("toml")
	viper.SetConfigName("config.toml")

	viper.SetDefault("db", path.Join(cfgDir, "mcp.db"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			cobra.CheckErr(err)
		}
	}

	dsnURI = viper.GetString("db")

	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		logger.Error("error connecting to database", "err", err, "uri", dsnURI)
		os.Exit(1)
	}
	defer db.Close()

	dir, err := iofs.New(dbDir, "db")
	cobra.CheckErr(err)
	defer dir.Close()

	instance, err := sqlite.WithInstance(db, &sqlite.Config{})
	cobra.CheckErr(err)

	migrations, err := migrate.NewWithInstance("iofs", dir, "sqlite:/"+dsnURI, instance)
	if err != nil {
		logger.Error("error creating migrations", "err", err, "uri", dsnURI)
		os.Exit(1)
	}
	defer migrations.Close()

	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("error running migrations", "err", err, "uri", dsnURI)
		os.Exit(1)
	}

	logger.Debug("migrations completed successfully", "uri", dsnURI)
}
