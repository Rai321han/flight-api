

package services

import "flight-api/models"

const flightIndex = "kibana_sample_data_flights"

// BuildSearchRequest constructs a full Elasticsearch Search API body
// from validated FlightFilters.
func BuildSearchRequest(f *models.FlightFilters) map[string]any {
    filter := []map[string]any{}

    // --- Term matches ---
    termFields := []struct {
        field string
        value *string
    }{
        {"Carrier", f.Carrier},
        {"FlightNum", f.FlightNum},
        {"DestCountry", f.DestCountry},
        {"OriginCountry", f.OriginCountry},
        {"DestCityName", f.DestCity},
        {"OriginCityName", f.OriginCity},
        {"DestAirportID", f.DestAirportID},
        {"OriginAirportID", f.OriginAirportID},
        {"DestAirport", f.DestAirport},
        {"OriginAirport", f.OriginAirport},
    }

    for _, tf := range termFields {
        if tf.value != nil {
            filter = append(filter, map[string]any{
                "term": map[string]any{tf.field: *tf.value},
            })
        }
    }

    // --- Cancelled bool ---
    if f.Cancelled != nil {
        filter = append(filter, map[string]any{
            "term": map[string]any{"Cancelled": *f.Cancelled},
        })
    }

    // --- Date range ---
    if f.DateFrom != nil || f.DateTo != nil {
        dateRange := map[string]any{}
        if f.DateFrom != nil {
            dateRange["gte"] = f.DateFrom.Format("2006-01-02")
        }
        if f.DateTo != nil {
            dateRange["lte"] = f.DateTo.Format("2006-01-02")
        }
        filter = append(filter, map[string]any{
            "range": map[string]any{"timestamp": dateRange},
        })
    }

    // --- Price range ---
    if f.PriceMin != nil || f.PriceMax != nil {
        priceRange := map[string]any{}
        if f.PriceMin != nil {
            priceRange["gte"] = *f.PriceMin
        }
        if f.PriceMax != nil {
            priceRange["lte"] = *f.PriceMax
        }
        filter = append(filter, map[string]any{
            "range": map[string]any{"AvgTicketPrice": priceRange},
        })
    }

    // --- Geo distance: destination ---
    if f.DestLat != nil && f.DestLon != nil {
        filter = append(filter, map[string]any{
            "geo_distance": map[string]any{
                "distance":     "100km",
                "DestLocation": map[string]any{"lat": *f.DestLat, "lon": *f.DestLon},
            },
        })
    }

    // --- Geo distance: origin ---
    if f.OriginLat != nil && f.OriginLon != nil {
        filter = append(filter, map[string]any{
            "geo_distance": map[string]any{
                "distance":       "100km",
                "OriginLocation": map[string]any{"lat": *f.OriginLat, "lon": *f.OriginLon},
            },
        })
    }

    // --- Assemble query ---
    query := map[string]any{}
    if len(filter) > 0 {
        query["bool"] = map[string]any{"filter": filter}
    } else {
        query["match_all"] = map[string]any{}
    }

    sortField := f.SortBy
    if sortField == "" {
        sortField = "timestamp"
    }
    order := f.Order
    if order == "" {
        order = "asc"
    }

    from := (f.Page - 1) * f.Limit
    if from < 0 {
        from = 0
    }

    return map[string]any{
        "query": query,
        "sort":  []map[string]any{{sortField: map[string]any{"order": order}}},
        "size":  f.Limit,
        "from":  from,
        "_source": []string{
            "FlightNum", "Carrier", "FlightDelayMin", "FlightDelayType",
            "OriginCityName", "OriginCountry", "OriginAirportID", "OriginAirport",
            "OriginLocation",
            "DestCityName", "DestCountry", "DestAirportID", "DestAirport",
            "DestLocation",
            "AvgTicketPrice", "Cancelled", "timestamp",
        },
    }
}

// IndexName returns the Elasticsearch index for flights.
func IndexName() string {
    return flightIndex
}
 