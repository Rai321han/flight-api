package services

import (
	"fmt"
	"log"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient     *elasticsearch.Client
	esClientErr  error
	esOnce       sync.Once
)

// GetESClient returns a lazily-initialised singleton Elasticsearch client.
// If initialisation previously failed the cached error is returned immediately.
func GetESClient() (*elasticsearch.Client, error) {
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
			esClientErr = fmt.Errorf("elasticsearch client init failed: %w", err)
			log.Printf("[ES] %v", esClientErr)
			return
		}

		esClient = client
		log.Printf("[ES] Client initialised → %s", host)
	})

	return esClient, esClientErr
}