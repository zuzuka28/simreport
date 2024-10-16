package config

import (
	"fmt"
	"simrep/api/amqp/asyncanalyze/consumer"
	"simrep/api/amqp/asyncanalyze/producer"
	"simrep/internal/repository/analyze"
	"simrep/internal/repository/document"
	"simrep/internal/repository/documentfile"
	"simrep/internal/repository/image"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"

	"github.com/ilyakaznacheev/cleanenv"
)

type StartupIndexOpts struct {
	Index         string `yaml:"index"`
	UpdateMapping bool   `yaml:"updateMapping"`
	CreateMapping bool   `yaml:"createMapping"`
	MappingPath   string `yaml:"mappingPath"`
	Alias         string `yaml:"alias"`
	Shards        int    `yaml:"shards"`
	Replics       int    `yaml:"replics"`
}

type Elastic struct {
	Hosts []string `json:"hosts"`
}

type Config struct {
	Port                 int                `yaml:"port"`
	S3                   minioutil.Config   `yaml:"s3"`
	Elastic              elasticutil.Config `yaml:"elastic"`
	VectorizerService    string             `yaml:"vectorizerService"`
	AnalyzeProducer      producer.Config    `yaml:"analyzeProducer"`
	AnalyzeConsumer      consumer.Config    `yaml:"analyzeConsumer"`
	ImageRepo            image.Opts         `yaml:"imageRepo"`
	DocumentFileRepo     documentfile.Opts  `yaml:"documentFileRepo"`
	DocumentRepo         document.Opts      `yaml:"documentRepo"`
	AnalyzedDocumentRepo analyze.Opts       `yaml:"analyzedDocumentRepo"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
