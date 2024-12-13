package main

import (
	"context"
	"errors"
	"mcp/internal/server"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
)

var (
	cmdServeStdio = &cobra.Command{
		Use:   "stdio",
		Short: "Start mcp as a stdio server.",
		Run: func(cmd *cobra.Command, args []string) {
			g, ctx := errgroup.WithContext(context.Background())

			g.Go(func() error {
				<-interrupts()
				logger.Info("interrupted")
				return context.Canceled
			})

			g.Go(func() error {
				_, err := server.NewServer(ctx, logger, os.Stdin, os.Stdout)
				if err != nil {
					return err
				}

				// Now that the server is set up, we want to load up configured clients and start them.
				// This will be done in a separate goroutine so that the server can start handling requests.

				logger.Info("server started")
				return nil
			})

			if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
				logger.Error("server error", "err", err)
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
