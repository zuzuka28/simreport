package minioutil

type Config struct {
	Endpoint        string   `yaml:"endpoint"`
	AccessKeyID     string   `yaml:"accessKeyID"`
	SecletAccessKey string   `yaml:"secletAccessKey"`
	Buckets         []string `yaml:"buckets"`
}
