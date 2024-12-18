package main

import (
	"context"
	"mcp/internal/integrations/sql"
	localbroker "mcp/internal/local_broker"
	docker_runner "mcp/internal/server_runner/docker"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cmdServeStdio = &cobra.Command{
		Use:   "stdio",
		Short: "Start mcp as a stdio server.",
		Run: func(cmd *cobra.Command, args []string) {
			g, ctx := errgroup.WithContext(cmd.Context())

			g.Go(func() error {
				<-interrupts()
				return context.Canceled
			})

			dsn := viper.GetString("db")
			logger.Debug("using database", "dsn", dsn)
			integRepo, err := sql.NewSQLDatabaseIntegrationsRepository(ctx, logger, dsn)
			if err != nil {
				logger.Error("error while creating integrations repository", "err", err)
				os.Exit(1)
			}
			defer integRepo.Close()

			logger.Debug("database up, starting docker runner")

			runner, err := docker_runner.NewDockerServerRunner(ctx, logger, docker_runner.DockerServerOptions{})
			if err != nil {
				logger.Error("error while creating docker server runner", "err", err)
				os.Exit(1)
			}
			defer runner.Close()

			g.Go(func() error {
				logger.Debug("docker runner up, starting local broker")

				broker := localbroker.NewLocalBroker(ctx, logger, integRepo, runner, os.Stdin, os.Stdout)
				defer broker.Close()

				if err := broker.Run(ctx); err != nil && err != context.Canceled {
					logger.Error("error while running local broker", "err", err)
					return err
				}

				logger.Debug("local broker finished")

				return nil
			})

			if err := g.Wait(); err != nil && err != context.Canceled {
				logger.Error("error while running server", "err", err)
				os.Exit(1)
			}
		},
	}
)

// interrupts returns a channel that is closed when an interrupt signal is received.
func interrupts() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		defer close(c)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
	}()
	return c
}
