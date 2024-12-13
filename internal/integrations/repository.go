package integrations

import (
	"database/sql"
	"mcp/internal/registry"
	"sync"
)

type InstalledIntegration struct {
	Id       string
	Manifest *registry.IntegrationManifest
}

type IntegrationsChangedEvent struct {
}

type IntegrationsChangedCallback func(e *IntegrationsChangedEvent)

type RemoveCallbackFunction func()

type IntegrationsRepository interface {
	InstallIntegration(m *registry.IntegrationManifest) (*InstalledIntegration, error)
	ListIntegrations() ([]*InstalledIntegration, error)
	UninstallIntegration(i *InstalledIntegration) error

	OnIntegrationsChanged(cb IntegrationsChangedCallback) RemoveCallbackFunction
}

var _ IntegrationsRepository = &databaseIntegrationsRepository{}

type databaseIntegrationsRepository struct {
	callbacks   map[struct{}]IntegrationsChangedCallback
	callbacksMu sync.RWMutex
}

func NewSQLDatabaseIntegrationsRepository(db *sql.DB) IntegrationsRepository {
	return &databaseIntegrationsRepository{
		callbacks: make(map[struct{}]IntegrationsChangedCallback),
	}
}

func (r *databaseIntegrationsRepository) InstallIntegration(m *registry.IntegrationManifest) (*InstalledIntegration, error) {
	return nil, nil
}

func (r *databaseIntegrationsRepository) ListIntegrations() ([]*InstalledIntegration, error) {
	return nil, nil
}

func (r *databaseIntegrationsRepository) UninstallIntegration(i *InstalledIntegration) error {
	return nil
}

func (r *databaseIntegrationsRepository) OnIntegrationsChanged(cb IntegrationsChangedCallback) RemoveCallbackFunction {
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
