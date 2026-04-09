package services

import (
	"fmt"
	"log"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient    *elasticsearch.Client
	esClientErr error
	esOnce      sync.Once
)

// GetESClient returns a lazily-initialised singleton Elasticsearch client.
// It reads connection settings from Beego app config on the first call and
// caches both the client and any initialisation error for all subsequent calls.
//
// Config keys:
//   - es.host     : full Elasticsearch base URL (default: http://localhost:9200)
//   - es.username : optional basic-auth username
//   - es.password : optional basic-auth password (ignored when username is absent)
//
// Returns:
//   - *elasticsearch.Client: ready-to-use client, or nil on failure
//   - error: non-nil when initialisation failed; the same error is returned on
//     every subsequent call — retrying will not recover because sync.Once
//     guarantees the constructor runs exactly once
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