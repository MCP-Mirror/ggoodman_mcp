package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/document"
	"github.com/blevesearch/bleve/v2/mapping"
	index "github.com/blevesearch/bleve_index_api"

	"github.com/google/uuid"
)

var _ RegistryClient = &fakeClient{}

type fakeClient struct {
	logger  *slog.Logger
	pkgs    []IntegrationManifest
	mapping mapping.IndexMapping
	index   bleve.Index
}

func NewFakeClient(logger *slog.Logger) (*fakeClient, error) {
	c := &fakeClient{
		logger: logger,
	}

	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = en.AnalyzerName

	pkgMapping := bleve.NewDocumentMapping()
	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Analyzer = "en"
	pkgMapping.AddFieldMappingsAt("Name", nameFieldMapping)
	descFieldMapping := bleve.NewTextFieldMapping()
	descFieldMapping.Analyzer = "en"
	pkgMapping.AddFieldMappingsAt("Description", descFieldMapping)
	mapping.AddDocumentMapping("package", pkgMapping)
	mapping.DefaultMapping = pkgMapping

	bIndex, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	c.index = bIndex
	c.mapping = mapping

	pkgs, err := loadFakePackages()
	if err != nil {
		return nil, err
	}

	batch := bIndex.NewBatch()

	for _, p := range pkgs {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		m := IntegrationSearchResult{
			Id:          id.String(),
			Name:        p.Name,
			Version:     "0.0.1",
			Description: p.Description,
		}

		doc := document.NewDocument(id.String())

		if err := mapping.MapDocument(doc, m); err != nil {
			return nil, err
		}

		mess, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("error marshalling document %s: %w", id.String(), err)
		}

		nameField := document.NewTextField("Name", nil, []byte(p.Name))
		descField := document.NewTextField("Description", nil, []byte(p.Name))
		sourceField := document.NewTextFieldWithIndexingOptions(
			"_source", nil, mess, index.StoreField)

		nd := doc.AddField(nameField).AddField(descField).AddField(sourceField)

		if err := batch.IndexAdvanced(nd); err != nil {
			return nil, err
		}

		c.pkgs = append(c.pkgs, IntegrationManifest{
			Id:          id.String(),
			Name:        p.Name,
			Version:     "0.0.1",
			Description: p.Description,
		})
	}

	if err := bIndex.Batch(batch); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *fakeClient) GetIntegrationManifestByNameAndVersion(name, version string) (*IntegrationManifest, error) {
	return nil, nil
}

func (c *fakeClient) SearchIntegrations(ctx context.Context, terms ...string) ([]*IntegrationSearchResult, error) {
	q := bleve.NewQueryStringQuery(strings.Join(terms, " "))
	sr := bleve.NewSearchRequest(q)
	sr.Fields = []string{"_source", "Name", "Description"}

	results, err := c.index.Search(sr)
	if err != nil {
		return nil, err
	}

	var pkgs []*IntegrationSearchResult

	for _, hit := range results.Hits {
		b, ok := hit.Fields["_source"]
		if !ok {
			return nil, fmt.Errorf("no _source field found for document %s", hit.ID)
		}

		bstr, ok := b.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", b)
		}

		var sr IntegrationSearchResult
		if err := json.Unmarshal([]byte(bstr), &sr); err != nil {
			return nil, fmt.Errorf("error unmarshalling document %s: %w", hit.ID, err)
		}

		pkgs = append(pkgs, &sr)
	}

	return pkgs, nil
}
