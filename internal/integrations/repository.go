package integrations

import (
	"database/sql"
	"mcp/internal/registry"
)

type InstalledIntegration struct {
	Id       string
	Manifest *registry.IntegrationManifest
}

type IntegrationsRepository interface {
	InstallIntegration(m *registry.IntegrationManifest) (*InstalledIntegration, error)
	ListIntegrations() ([]*InstalledIntegration, error)
	UninstallIntegration(i *InstalledIntegration) error
}

var _ IntegrationsRepository = &databaseIntegrationsRepository{}

type databaseIntegrationsRepository struct {
}

func NewSQLDatabaseIntegrationsRepository(db *sql.DB) IntegrationsRepository {
	return &databaseIntegrationsRepository{}
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
