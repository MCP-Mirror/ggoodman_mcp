package registry

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
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
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	c.index = index
	c.mapping = mapping

	pkgs, err := loadFakePackages()
	if err != nil {
		return nil, err
	}

	batch := index.NewBatch()

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

		var mess bytes.Buffer

		enc := gob.NewEncoder(&mess)
		enc.Encode(m)

		field := document.NewTextFieldWithIndexingOptions(
			"_source", nil, mess.Bytes(), document.StoreField)
		nd := doc.AddField(field)

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

	if err := index.Batch(batch); err != nil {
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
	sr.Fields = []string{"_source"}

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

		messout := bytes.NewBuffer([]byte(fmt.Sprintf("%v", b)))
		dec := gob.NewDecoder(messout)

		var sr IntegrationSearchResult
		err = dec.Decode(&sr)
		if err != nil {
			return nil, fmt.Errorf("error decoding document %s: %w", hit.ID, err)
		}

		pkgs = append(pkgs, &sr)
	}

	return pkgs, nil
}
