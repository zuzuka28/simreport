package config

import (
	"fmt"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"
	documentcmd "github.com/zuzuka28/simreport/prj/document/pkg/cmd"
	fulltextindexcmd "github.com/zuzuka28/simreport/prj/fulltextindex/pkg/cmd"
	semanticindexcmd "github.com/zuzuka28/simreport/prj/semanticindex/pkg/cmd"
	shingleindexcmd "github.com/zuzuka28/simreport/prj/shingleindex/pkg/cmd"
	similaritycmd "github.com/zuzuka28/simreport/prj/similarity/pkg/cmd"
)

type Config struct {
	Port        int `yaml:"port"`
	MetricsPort int `yaml:"metricsPort"`

	Nats    string             `yaml:"nats"`
	Elastic elasticutil.Config `yaml:"elastic"`
	S3      minioutil.Config   `yaml:"s3"`

	DocumentService   documentcmd.Config   `yaml:"documentService"`
	SimilarityService similaritycmd.Config `yaml:"similarityService"`

	ShingleIndexService  shingleindexcmd.Config  `yaml:"shingleIndexService"`
	FulltextIndexService fulltextindexcmd.Config `yaml:"fulltextIndexService"`
	SemanticIndexService semanticindexcmd.Config `yaml:"semanticIndexService"`
}

func New(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if cfg.MetricsPort == 0 {
		cfg.MetricsPort = 9000
	}

	documentCfg, err := getDocumentServiceConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("get document service config: %w", err)
	}

	cfg.DocumentService = documentCfg

	similarityCfg, err := getSimilarityServiceConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("get similarity service config: %w", err)
	}

	cfg.SimilarityService = similarityCfg

	shingleCfg, err := getShingleIndexServiceConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("get shingleindex service config: %w", err)
	}

	cfg.ShingleIndexService = shingleCfg

	fulltextCfg, err := getFulltextIndexServiceConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("get fulltextindex service config: %w", err)
	}

	cfg.FulltextIndexService = fulltextCfg

	semanticCfg, err := getSemanticIndexServiceConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("get semanticindex service config: %w", err)
	}

	cfg.SemanticIndexService = semanticCfg

	return &cfg, nil
}

func getDocumentServiceConfig(cfg Config) (documentcmd.Config, error) {
	result := cfg.DocumentService

	result.Nats = cfg.Nats
	result.Elastic = cfg.Elastic
	result.S3 = cfg.S3

	port := result.Port

	if port == 0 {
		apiPort, err := getFreePort()
		if err != nil {
			return documentcmd.Config{}, fmt.Errorf("get free port for API: %w", err)
		}

		port = apiPort
	}

	result.Port = port

	return result, nil
}

func getSimilarityServiceConfig(cfg Config) (similaritycmd.Config, error) {
	result := cfg.SimilarityService

	result.Nats = cfg.Nats
	result.Elastic = cfg.Elastic
	result.S3 = cfg.S3

	port := result.Port

	if port == 0 {
		apiPort, err := getFreePort()
		if err != nil {
			return similaritycmd.Config{}, fmt.Errorf("get free port for API: %w", err)
		}

		port = apiPort
	}

	result.Port = port

	return result, nil
}

func getShingleIndexServiceConfig(cfg Config) (shingleindexcmd.Config, error) {
	result := cfg.ShingleIndexService

	result.Nats = cfg.Nats
	result.Elastic = cfg.Elastic
	result.S3 = cfg.S3

	return result, nil
}

func getFulltextIndexServiceConfig(cfg Config) (fulltextindexcmd.Config, error) {
	result := cfg.FulltextIndexService

	result.Nats = cfg.Nats
	result.Elastic = cfg.Elastic
	result.S3 = cfg.S3

	return result, nil
}

func getSemanticIndexServiceConfig(cfg Config) (semanticindexcmd.Config, error) {
	result := cfg.SemanticIndexService

	result.Nats = cfg.Nats
	result.Elastic = cfg.Elastic
	result.S3 = cfg.S3

	return result, nil
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer func() { _ = l.Close() }()

	return l.Addr().(*net.TCPAddr).Port, nil
}
