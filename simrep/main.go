package main

import (
	"fmt"
	"simrep/api/amqp/asyncanalyze/consumer"
	"simrep/api/amqp/asyncanalyze/producer"
	"simrep/internal/config"
	"simrep/internal/repository/analyze"
	"simrep/internal/repository/document"
	"simrep/internal/repository/documentfile"
	"simrep/internal/repository/image"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"

	"gopkg.in/yaml.v3"
)

func main() {
	c := config.Config{
		Port:              0,
		S3:                minioutil.Config{Endpoint: "", AccessKeyID: "", SecletAccessKey: "", Buckets: []string{}},
		Elastic:           elasticutil.Config{Hosts: []string{}, IndexOpts: []elasticutil.StartupIndexConfig{{Index: "", UpdateMapping: false, CreateMapping: false, MappingPath: "", Alias: "", Shards: 0, Replics: 0}}},
		VectorizerService: "",
		AnalyzeProducer: producer.Config{
			DSN:          "",
			ExchangeName: "",
			QueueName:    "",
			RoutingKey:   "",
			MaxRetries:   0,
		},
		AnalyzeConsumer: consumer.Config{
			DSN:        "",
			QueueName:  "",
			MaxRetries: 0,
		},
		ImageRepo:        image.Opts{Bucket: ""},
		DocumentFileRepo: documentfile.Opts{Bucket: ""},
		DocumentRepo:     document.Opts{Index: ""},
		AnalyzedDocumentRepo: analyze.Opts{
			Index: "",
		},
	}

	out, _ := yaml.Marshal(c)

	fmt.Println(string(out))
}
