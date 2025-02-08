package elasticutil

import (
	"net/http"
)

type StartupIndexConfig struct {
	Index         string `yaml:"index"`
	UpdateMapping bool   `yaml:"updateMapping"`
	CreateMapping bool   `yaml:"createMapping"`
	MappingPath   string `yaml:"mappingPath"`
	Alias         string `yaml:"alias"`
	Shards        int    `yaml:"shards"`
	Replics       int    `yaml:"replics"`
}

type Config struct {
	Hosts     []string             `yaml:"hosts"`
	IndexOpts []StartupIndexConfig `yaml:"startupIndex"`

	Transport http.RoundTripper `yaml:"-"`
}
