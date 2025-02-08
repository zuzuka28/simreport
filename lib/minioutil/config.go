package minioutil

import "net/http"

type Config struct {
	Endpoint        string   `yaml:"endpoint"`
	AccessKeyID     string   `yaml:"accessKeyID"`
	SecletAccessKey string   `yaml:"secletAccessKey"`
	Buckets         []string `yaml:"buckets"`

	Transport http.RoundTripper `yaml:"-"`
}
