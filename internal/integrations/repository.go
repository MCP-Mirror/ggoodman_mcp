package integrations

import (
	"context"
	"mcp/internal/registry"
)

type InstalledIntegration struct {
	Id       string
	Manifest *registry.IntegrationManifest
}

type IntegrationsChangedEventType int

const (
	IntegrationsChangedEventTypeAdded IntegrationsChangedEventType = iota
	IntegrationsChangedEventTypeRemoved
)

type IntegrationsChangedEvent struct {
	Type        IntegrationsChangedEventType
	Integration InstalledIntegration
}

type IntegrationsChangedCallback func(e *IntegrationsChangedEvent)

type HandlerRemover interface {
	Close()
}

type IntegrationsRepository interface {
	Close() error

	InstallIntegration(ctx context.Context, m *registry.IntegrationManifest) (*InstalledIntegration, error)
	ListIntegrations(ctx context.Context) ([]*InstalledIntegration, error)
	UninstallIntegration(ctx context.Context, i *InstalledIntegration) error

	OnIntegrationsChanged(cb IntegrationsChangedCallback) HandlerRemover
}
