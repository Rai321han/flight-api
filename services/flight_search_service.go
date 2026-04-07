package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"flight-api/models"
)

// FlightSearchService is the live Elasticsearch implementation of models.FlightService.
type FlightSearchService struct{}

// NewFlightSearchService constructs a FlightSearchService.
func NewFlightSearchService() *FlightSearchService {
	return &FlightSearchService{}
}

// SearchFlights executes the ES query and returns source documents,
// total hit count, and any error.
func (s *FlightSearchService) SearchFlights(
	filters *models.FlightFilters,
) ([]map[string]any, int, error) {

	client, err := GetESClient()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ES client: %w", err)
	}

	bodyBytes, err := json.Marshal(BuildSearchRequest(filters))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialise query: %w", err)
	}

	log.Printf("[FlightSearchService] query → %s", string(bodyBytes))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(IndexName()),
		client.Search.WithBody(bytes.NewReader(bodyBytes)),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("elasticsearch search error: %w", err)
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read ES response body: %w", err)
	}

	if res.IsError() {
		return nil, 0, fmt.Errorf("ES returned HTTP %s: %s", res.Status(), string(rawBody))
	}

	var esResp esSearchResponse
	if err := json.Unmarshal(rawBody, &esResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode ES response: %w", err)
	}

	results := make([]map[string]any, 0, len(esResp.Hits.Hits))
	for _, hit := range esResp.Hits.Hits {
		results = append(results, hit.Source)
	}

	total := esResp.Hits.Total.Value
	log.Printf("[FlightSearchService] returned %d / %d hits", len(results), total)

	return results, total, nil
}

// --- ES response envelope ---

type esSearchResponse struct {
	Hits esHitsWrapper `json:"hits"`
}

type esHitsWrapper struct {
	Total esTotal `json:"total"`
	Hits  []esHit `json:"hits"`
}

type esTotal struct {
	Value int `json:"value"`
}

type esHit struct {
	Source map[string]any `json:"_source"`
}