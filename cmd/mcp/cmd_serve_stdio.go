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
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			g, ctx := errgroup.WithContext(cmd.Context())

			g.Go(func() error {
				select {
				case <-ctx.Done():
					return nil
				case <-interrupts():
					return nil
				}
			})

			g.Go(func() error {
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

				logger.Debug("docker runner up, starting local broker")

				broker := localbroker.NewLocalBroker(ctx, logger, integRepo, runner, os.Stdin, os.Stdout)
				defer broker.Close()

				if err := broker.Run(ctx); err != nil {
					if err != localbroker.ErrConnectionClosed {
						logger.Error("error while running local broker", "err", err)
					}

					return err
				}

				logger.Debug("local broker finished")

				return nil
			})

			if err := g.Wait(); err != nil && err != context.Canceled && err != localbroker.ErrConnectionClosed {
				logger.Error("error while running server", "err", err)
				os.Exit(1)
			}

			os.Exit(0)
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
