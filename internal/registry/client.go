package registry

import "context"

type IntegrationSearchResult struct {
	Id          string
	Name        string
	Version     string
	Description string
}

type IntegrationManifest struct {
	Id          string
	Name        string
	Version     string
	Description string
	Vendor      string
	SourceURL   string
	License     string
	Homepage    string
	Runtime     string
}

type RegistryClient interface {
	GetIntegrationManifestByNameAndVersion(name, version string) (*IntegrationManifest, error)
	SearchIntegrations(ctx context.Context, terms ...string) ([]*IntegrationSearchResult, error)
}
