package services

import (
	"fmt"
	"log"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient *elasticsearch.Client
	esOnce   sync.Once
)

// GetESClient returns a lazily-initialised singleton Elasticsearch client.
func GetESClient() (*elasticsearch.Client, error) {
	var initErr error

	esOnce.Do(func() {
		host, _ := beego.AppConfig.String("es.host")
		if host == "" {
			host = "http://localhost:9200"
		}

		username, _ := beego.AppConfig.String("es.username")
		password, _ := beego.AppConfig.String("es.password")

		cfg := elasticsearch.Config{
			Addresses: []string{host},
		}
		if username != "" && password != "" {
			cfg.Username = username
			cfg.Password = password
		}

		client, err := elasticsearch.NewClient(cfg)
		if err != nil {
			initErr = fmt.Errorf("elasticsearch client init failed: %w", err)
			log.Printf("[ES] %v", initErr)
			return
		}

		esClient = client
		log.Printf("[ES] Client initialised → %s", host)
	})

	if initErr != nil {
		return nil, initErr
	}
	return esClient, nil
}