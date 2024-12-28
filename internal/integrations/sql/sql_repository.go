package sql

import (
	"context"
	"database/sql"
	"embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"mcp/internal/integrations"
	"mcp/internal/integrations/sql/internal"
	"mcp/internal/registry"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsDir embed.FS

var _ integrations.IntegrationsRepository = &databaseIntegrationsRepository{}

type databaseIntegrationsRepository struct {
	callbacks   map[struct{}]integrations.IntegrationsChangedCallback
	callbacksMu sync.RWMutex

	logger *slog.Logger
	db     *sql.DB
}

func NewSQLDatabaseIntegrationsRepository(ctx context.Context, logger *slog.Logger, dsnURI string) (integrations.IntegrationsRepository, error) {
	var err error

	localDb, err := sql.Open("sqlite", "file://"+dsnURI+"?_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	defer localDb.Close()

	dir, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error creating migrations source: %w", err)
	}
	defer dir.Close()

	instance, err := sqlite.WithInstance(localDb, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating migrations instance: %w", err)
	}

	migrations, err := migrate.NewWithInstance("iofs", dir, "sqlite:/"+dsnURI+"?_pragma=journal_mode(WAL)", instance)
	if err != nil {
		return nil, fmt.Errorf("error creating migrations: %w", err)
	}
	defer migrations.Close()

	err = migrations.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return nil, fmt.Errorf("error running migrations: %w", err)
		}
		err = nil
	}
	if err != migrate.ErrNoChange {
		logger.Debug("migrations completed successfully", "uri", dsnURI)
	}

	db, err := sql.Open("sqlite", "file://"+dsnURI+"?_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &databaseIntegrationsRepository{
		callbacks: make(map[struct{}]integrations.IntegrationsChangedCallback),
		db:        db,
		logger:    logger,
	}, nil
}

func (r *databaseIntegrationsRepository) Close() error {
	r.logger.Debug("closing database")
	return r.db.Close()
}

func (r *databaseIntegrationsRepository) InstallIntegration(ctx context.Context, m *registry.IntegrationManifest) (*integrations.InstalledIntegration, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbtx := internal.New(r.db)

	i, err := dbtx.CreateIntegration(ctx, internal.CreateIntegrationParams{
		Name:         m.Name,
		Description:  m.Description,
		Vendor:       m.Vendor,
		SourceUrl:    m.SourceURL,
		Homepage:     m.Homepage,
		License:      m.License,
		Instructions: []byte{},
	})
	if err != nil {
		return nil, fmt.Errorf("error installing integration: %w", err)
	}

	return &integrations.InstalledIntegration{
		Id:       base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", i.ID))),
		Manifest: m,
	}, nil
}

func (r *databaseIntegrationsRepository) ListIntegrations(ctx context.Context) ([]*integrations.InstalledIntegration, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbtx := internal.New(r.db)

	listed, err := dbtx.ListIntegrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing integrations: %w", err)
	}

	var installed []*integrations.InstalledIntegration
	for _, i := range listed {
		installed = append(installed, &integrations.InstalledIntegration{
			Id: base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", i.ID))),
			Manifest: &registry.IntegrationManifest{
				Name:        i.Name,
				Description: i.Description,
				Vendor:      i.Vendor,
				SourceURL:   i.SourceUrl,
				Homepage:    i.Homepage,
				License:     i.License,
				// TODO: Remove hard-coded values
				Runtime: "node",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-github"},
			},
		})
	}

	return installed, nil
}

func (r *databaseIntegrationsRepository) UninstallIntegration(ctx context.Context, i *integrations.InstalledIntegration) error {
	return nil
}

func (r *databaseIntegrationsRepository) OnIntegrationsChanged(cb integrations.IntegrationsChangedCallback) integrations.HandlerRemover {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()

	key := struct{}{}
	r.callbacks[key] = cb

	return &handlerRemover[integrations.IntegrationsChangedCallback]{
		callbacks: r.callbacks,
		mu:        &r.callbacksMu,
		key:       key,
	}
}

type handlerRemover[T any] struct {
	mu        *sync.RWMutex
	key       struct{}
	callbacks map[struct{}]T
}

func (hr *handlerRemover[T]) Close() {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	delete(hr.callbacks, hr.key)
}
