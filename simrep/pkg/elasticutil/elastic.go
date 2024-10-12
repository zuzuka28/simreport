package elasticutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func NewClientWithStartup(
	ctx context.Context,
	cfg Config,
) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{ //nolint:exhaustruct
		Addresses: cfg.Hosts,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("creating ElasticSearch client: %w", err)
	}

	for _, startupCfg := range cfg.IndexOpts {
		if err := startupIndex(ctx, client, startupCfg); err != nil {
			return nil, fmt.Errorf("startup index %s: %w", startupCfg.Index, err)
		}
	}

	return client, nil
}

func startupIndex(
	ctx context.Context,
	client *elasticsearch.Client,
	indexCfg StartupIndexConfig,
) error {
	indexExists, err := checkIndexExists(ctx, client, indexCfg.Index)
	if err != nil {
		return err
	}

	if !indexExists && indexCfg.CreateMapping {
		if err := createIndex(ctx, client, indexCfg); err != nil {
			return err
		}
	}

	if indexCfg.UpdateMapping {
		if err := updateIndexMapping(ctx, client, indexCfg); err != nil {
			return err
		}
	}

	if indexCfg.Alias != "" {
		if err := createAlias(ctx, client, indexCfg.Index, indexCfg.Alias); err != nil {
			return err
		}
	}

	return nil
}

func checkIndexExists(
	ctx context.Context,
	client *elasticsearch.Client,
	indexName string,
) (bool, error) {
	res, err := client.Indices.Exists(
		[]string{indexName},
		client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("check if index exists: %w", err)
	}

	defer res.Body.Close()

	return res.StatusCode != http.StatusNotFound, nil
}

func createIndex(
	ctx context.Context,
	client *elasticsearch.Client,
	indexCfg StartupIndexConfig,
) error {
	if indexCfg.MappingPath == "" {
		return nil
	}

	mapping, err := loadMapping(indexCfg.MappingPath)
	if err != nil {
		return fmt.Errorf("load mapping: %w", err)
	}

	mapping.Settings = map[string]any{
		"number_of_shards":   indexCfg.Shards,
		"number_of_replicas": indexCfg.Replics,
	}

	m, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("marshal mapping: %w", err)
	}

	res, err := client.Indices.Create(
		indexCfg.Index,
		client.Indices.Create.WithBody(bytes.NewReader(m)),
		client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	defer res.Body.Close()

	if err := IsErr(res); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	slog.Error("Index created", "index", indexCfg.Index)

	return nil
}

func updateIndexMapping(
	ctx context.Context,
	client *elasticsearch.Client,
	indexCfg StartupIndexConfig,
) error {
	if indexCfg.MappingPath == "" {
		return nil
	}

	body, err := loadMapping(indexCfg.MappingPath)
	if err != nil {
		return fmt.Errorf("load mapping for update: %w", err)
	}

	body.Settings = nil

	m, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal mapping: %w", err)
	}

	res, err := client.Indices.PutMapping(
		[]string{indexCfg.Index},
		bytes.NewReader(m),
		client.Indices.PutMapping.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("update mapping: %w", err)
	}

	defer res.Body.Close()

	if err := IsErr(res); err != nil {
		return fmt.Errorf("update mapping: %w", err)
	}

	slog.Error("Mapping updated", "index", indexCfg.Index)

	return nil
}

func createAlias(
	ctx context.Context,
	client *elasticsearch.Client,
	index, alias string,
) error {
	body := fmt.Sprintf(`{
		"actions": [
			{"add": {"index": "%s", "alias": "%s"}}
		]
	}`, index, alias)

	res, err := client.Indices.UpdateAliases(
		bytes.NewReader([]byte(body)),
		client.Indices.UpdateAliases.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("create alias: %w", err)
	}

	defer res.Body.Close()

	if err := IsErr(res); err != nil {
		return fmt.Errorf("creaet alias: %w", err)
	}

	slog.Info("Alias created", "alias", alias, "index", index)

	return nil
}

func loadMapping(mappingPath string) (*indexConfig, error) {
	file, err := os.Open(mappingPath)
	if err != nil {
		return nil, fmt.Errorf("open mapping file: %w", err)
	}

	defer file.Close()

	var mapping indexConfig

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read mapping file: %w", err)
	}

	if err := json.Unmarshal(raw, &mapping); err != nil {
		return nil, fmt.Errorf("unmarshal mapping: %w", err)
	}

	return &mapping, nil
}
