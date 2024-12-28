package main

import (
	"log/slog"
	"os"
	"path"

	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

var (
	logger *slog.Logger

	logLevelArg string
	logLevel    = &slog.LevelVar{}

	cmdRoot = &cobra.Command{
		Use:   "mcp",
		Short: "MCP is a cli tool for managing all of your model context needs via the Model Context Protocol.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	cmdRoot.PersistentFlags().StringVar(&logLevelArg, "log-level", "info", "log level among \"debug\", \"info\" or \"error\"")

	cmdRoot.AddCommand(cmdPackage)
	cmdRoot.AddCommand(cmdRegistry)
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

	viper.SetDefault("logfile", path.Join(cfgDir, "debug.log"))
	viper.SetDefault("db", path.Join(cfgDir, "db.sqlite"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			cobra.CheckErr(err)
		}
	}

	logFile := viper.GetString("logfile")
	logFd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	cobra.CheckErr(err)

	logger = slog.New(slogmulti.Fanout(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}), slog.NewTextHandler(logFd, &slog.HandlerOptions{
		Level: logLevel,
	})))
}
