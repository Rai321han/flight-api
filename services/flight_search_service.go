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

// FlightSearchService is the live Elasticsearch implementation of
// models.FlightService. It holds no state of its own — all configuration
// is read from the singleton ES client and the filters passed per request.
type FlightSearchService struct{}

// NewFlightSearchService constructs and returns a FlightSearchService.
// Elasticsearch client initialisation is deferred to the first search call
// via GetESClient, so this constructor never fails.
//
// Returns:
//   - *FlightSearchService: ready-to-use service instance
func NewFlightSearchService() *FlightSearchService {
	return &FlightSearchService{}
}

// SearchFlights executes an Elasticsearch query against the flights index
// using the criteria encoded in filters and returns the matching documents.
//
// Parameters:
//   - filters: validated FlightFilters containing all active search criteria,
//     sort preference, and pagination settings
//
// Returns:
//   - []map[string]any: source documents for the current page, each field
//     limited to the projection defined in flightSourceFields
//   - int: total number of documents matching the query before pagination
//   - error: *models.AppError on failure, nil on success
//
// Errors:
//   - 503: Elasticsearch client failed to initialise
//   - 500: search request body could not be serialised
//   - 502: network failure, unreadable response body, HTTP-level ES error,
//     or JSON decode failure from the ES response
func (s *FlightSearchService) SearchFlights(
	filters *models.FlightFilters,
) ([]map[string]any, int, error) {

	client, err := GetESClient()
	if err != nil {
		return nil, 0, models.NewAppError(503, "search backend unavailable", err)
	}

	bodyBytes, err := json.Marshal(BuildSearchRequest(filters))
	if err != nil {
		return nil, 0, models.NewAppError(500, "failed to build search query", err)
	}

	log.Printf("[FlightSearchService] executing query → %s", string(bodyBytes))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(IndexName()),
		client.Search.WithBody(bytes.NewReader(bodyBytes)),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, models.NewAppError(502, "search backend request failed", err)
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, models.NewAppError(502, "failed to read search backend response", err)
	}

	// Body is fully consumed before IsError so the connection is drained
	// regardless of outcome and the raw body can be included in the error
	// message for diagnostics.
	if res.IsError() {
		return nil, 0, models.NewAppError(
			502,
			"search backend returned an error",
			fmt.Errorf("ES status %s: %s", res.Status(), string(rawBody)),
		)
	}

	var esResp esSearchResponse
	if err := json.Unmarshal(rawBody, &esResp); err != nil {
		return nil, 0, models.NewAppError(502, "failed to decode search backend response", err)
	}

	results := make([]map[string]any, 0, len(esResp.Hits.Hits))
	for _, hit := range esResp.Hits.Hits {
		results = append(results, hit.Source)
	}

	total := esResp.Hits.Total.Value
	log.Printf("[FlightSearchService] page=%d limit=%d returned=%d total=%d",
		filters.Page, filters.Limit, len(results), total)

	return results, total, nil
}

// esSearchResponse is the top-level envelope returned by the Elasticsearch
// Search API. Only the fields consumed by SearchFlights are mapped; score
// and metadata fields are intentionally omitted.
type esSearchResponse struct {
	Hits esHitsWrapper `json:"hits"`
}

// esHitsWrapper wraps the hits array and the total hit count returned by ES.
type esHitsWrapper struct {
	Total esTotal `json:"total"`
	Hits  []esHit `json:"hits"`
}

// esTotal carries the total number of documents that matched the query.
// Value is an exact count when track_total_hits is true.
type esTotal struct {
	Value int `json:"value"`
}

// esHit represents a single document in the hits array.
// Only _source is extracted; score and metadata fields are ignored.
type esHit struct {
	Source map[string]any `json:"_source"`
}