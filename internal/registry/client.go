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
}

type RegistryClient interface {
	GetIntegrationManifestByNameAndVersion(name, version string) (*IntegrationManifest, error)
	SearchIntegrations(ctx context.Context, terms ...string) ([]*IntegrationSearchResult, error)
}

var _ RegistryClient = &FakeClient{}

type FakeClient struct {
}

func (c *FakeClient) GetIntegrationManifestByNameAndVersion(name, version string) (*IntegrationManifest, error) {
	return nil, nil
}

func (c *FakeClient) SearchIntegrations(ctx context.Context, terms ...string) ([]*IntegrationSearchResult, error) {
	return nil, nil
}
