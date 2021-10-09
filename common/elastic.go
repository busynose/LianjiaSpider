package common

import (
	"github.com/elastic/go-elasticsearch/v7"
)

var EsClient *elasticsearch.Client

func InitElastic() *elasticsearch.Client {
	var err error
	EsClient, err = elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	return EsClient
}
