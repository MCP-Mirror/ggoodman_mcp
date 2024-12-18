package sql

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"mcp/internal/integrations"
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

	localDb, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	defer func() {
		if err != nil {
			localDb.Close()
		}
	}()

	dir, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error creating migrations source: %w", err)
	}
	defer dir.Close()

	instance, err := sqlite.WithInstance(localDb, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating migrations instance: %w", err)
	}

	migrations, err := migrate.NewWithInstance("iofs", dir, "sqlite:/"+dsnURI, instance)
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

	db, err := sql.Open("sqlite", dsnURI)
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
	return nil, nil
}

var queryInstalledIntegrations = `
SELECT id, name, description, vendor, source_url, homepage, license, runtime
FROM integrations
`

func (r *databaseIntegrationsRepository) ListIntegrations(ctx context.Context) ([]*integrations.InstalledIntegration, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stmt, err := r.db.PrepareContext(ctx, queryInstalledIntegrations)
	if err != nil {
		return nil, fmt.Errorf("error preparing list integrations query: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error querying installed integrations: %w", err)
	}
	defer rows.Close()

	var installed []*integrations.InstalledIntegration

	for rows.Next() {
		var i integrations.InstalledIntegration

		if err := rows.Scan(&i.Id, &i.Manifest.Name, &i.Manifest.Description, &i.Manifest.Vendor, &i.Manifest.SourceURL, &i.Manifest.Homepage, &i.Manifest.License, &i.Manifest.Runtime); err != nil {
			return nil, fmt.Errorf("error scanning installed integration: %w", err)
		}

		installed = append(installed, &i)
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
