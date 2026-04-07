package models

// AggregationResult holds all aggregation buckets returned by the service.
type AggregationResult struct {
	FlightsPerCarrier      []BucketCount   `json:"flightsPerCarrier"`
	FlightsPerDestCountry  []BucketCount   `json:"flightsPerDestCountry"`
	FlightsPerOriginCountry []BucketCount  `json:"flightsPerOriginCountry"`
	AvgPricePerCarrier     []BucketAvg     `json:"avgPricePerCarrier"`
	CancelledPerCarrier    []BucketCount   `json:"cancelledPerCarrier"`
	DelayTypeBreakdown     []BucketCount   `json:"delayTypeBreakdown"`
	TopRoutes              []BucketCount   `json:"topRoutes"`
	PriceHistogram         []BucketCount   `json:"priceHistogram"`
	FlightsPerDay          []BucketCount   `json:"flightsPerDay"`
	OverallAvgPrice        float64         `json:"overallAvgPrice"`
}

// BucketCount represents a terms/histogram bucket with a doc count.
type BucketCount struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// BucketAvg represents a terms bucket with a nested avg metric.
type BucketAvg struct {
	Key   string  `json:"key"`
	Count int     `json:"count"`
	Avg   float64 `json:"avg"`
}