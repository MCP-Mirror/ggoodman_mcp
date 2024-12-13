package integrations

import (
	"context"
	"mcp/internal/registry"
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
	InstallIntegration(ctx context.Context, m *registry.IntegrationManifest) (*InstalledIntegration, error)
	ListIntegrations(ctx context.Context) ([]*InstalledIntegration, error)
	UninstallIntegration(ctx context.Context, i *InstalledIntegration) error

	OnIntegrationsChanged(cb IntegrationsChangedCallback) RemoveCallbackFunction
}
