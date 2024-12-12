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

// var (
// 	stdin  io.ReadCloser
// 	stdout io.WriteCloser
// )

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Starts the MCP server.",
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

// func init() {
// 	// stdin = newBufferedReadCloser(os.Stdin)
// 	// stdout = newBufferedWriteCloser(os.Stdout)
// }

// type bufferedReadCloser struct {
// 	*bufio.Reader
// 	closer io.Closer
// }

// func (brc *bufferedReadCloser) Close() error {
// 	return brc.closer.Close()
// }

// func newBufferedReadCloser(r io.ReadCloser) *bufferedReadCloser {
// 	return &bufferedReadCloser{
// 		Reader: bufio.NewReader(r),
// 		closer: r,
// 	}
// }

// type bufferedWriteCloser struct {
// 	*bufio.Writer
// 	closer io.Closer
// }

// func (brc *bufferedWriteCloser) Close() error {
// 	return brc.closer.Close()
// }

// func newBufferedWriteCloser(r io.WriteCloser) *bufferedWriteCloser {
// 	return &bufferedWriteCloser{
// 		Writer: bufio.NewWriter(r),
// 		closer: r,
// 	}
// }

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
