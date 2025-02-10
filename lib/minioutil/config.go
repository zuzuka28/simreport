package minioutil

import "net/http"

type Config struct {
	Endpoint        string   `yaml:"endpoint"`
	AccessKeyID     string   `yaml:"accessKeyID"`
	SecretAccessKey string   `yaml:"secretAccessKey"`
	Buckets         []string `yaml:"buckets"`

	Transport http.RoundTripper `yaml:"-"`
}
