package main

import (
	"context"
	"mcp/internal/integrations/sql"
	"mcp/internal/registry"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	cmdPackageInstall = &cobra.Command{
		Use:     "install <package[@version]>",
		Short:   "Install a package from the registry.",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			var manifest = registry.IntegrationManifest{
				Id:          "fake123",
				Name:        "GitHub MCP Server",
				Version:     "0.6.2",
				Description: "MCP server for using the GitHub API",
				Vendor:      "Anthropic, PBC (https://anthropic.com)",
				SourceURL:   "https://github.com/modelcontextprotocol/servers/tree/main/src/github",
				License:     "MIT",
				Homepage:    "https://modelcontextprotocol.io",
				Runtime:     "node",
				Command:     "npx",
				Args:        []string{"-y", "@modelcontextprotocol/server-github"},
				EnvVars: []registry.IntegrationEnvConfig{
					{
						Name:        "GITHUB_PERSONAL_ACCESS_TOKEN",
						Description: "A GitHub personal access token with the necessary permissions to access the GitHub API.",
						Default:     nil,
					},
				},
			}

			g, ctx := errgroup.WithContext(cmd.Context())

			go func() {

				<-interrupts()

				cancel()
			}()

			g.Go(func() error {
				dsn := viper.GetString("db")
				integRepo, err := sql.NewSQLDatabaseIntegrationsRepository(ctx, logger, dsn)
				if err != nil {
					logger.Error("error while creating integrations repository", "err", err)
					return err
				}
				defer integRepo.Close()

				if _, err = integRepo.InstallIntegration(ctx, &manifest); err != nil {
					logger.Info("Integration installation failed", "err", err)
					return err
				}

				logger.Info("Integration installed successfully", "name", manifest.Name, "version", manifest.Version)

				return nil
			})

			if err := g.Wait(); err != nil && err != context.Canceled {
				logger.Error("error while running server", "err", err)
				os.Exit(1)
			}

			os.Exit(0)
		},
	}
)
