package sql

import (
	"context"
	"database/sql"
	"mcp/internal/integrations"
	"mcp/internal/registry"
	"sync"
)

var _ integrations.IntegrationsRepository = &databaseIntegrationsRepository{}

type databaseIntegrationsRepository struct {
	callbacks   map[struct{}]integrations.IntegrationsChangedCallback
	callbacksMu sync.RWMutex
}

func NewSQLDatabaseIntegrationsRepository(db *sql.DB) (integrations.IntegrationsRepository, error) {
	return &databaseIntegrationsRepository{
		callbacks: make(map[struct{}]integrations.IntegrationsChangedCallback),
	}, nil
}

func (r *databaseIntegrationsRepository) InstallIntegration(ctx context.Context, m *registry.IntegrationManifest) (*integrations.InstalledIntegration, error) {
	return nil, nil
}

func (r *databaseIntegrationsRepository) ListIntegrations(ctx context.Context) ([]*integrations.InstalledIntegration, error) {
	return nil, nil
}

func (r *databaseIntegrationsRepository) UninstallIntegration(ctx context.Context, i *integrations.InstalledIntegration) error {
	return nil
}

func (r *databaseIntegrationsRepository) OnIntegrationsChanged(cb integrations.IntegrationsChangedCallback) integrations.RemoveCallbackFunction {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()

	key := struct{}{}
	r.callbacks[key] = cb

	return func() {
		r.callbacksMu.Lock()
		defer r.callbacksMu.Unlock()
		delete(r.callbacks, key)
	}
}
