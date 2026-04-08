# Flight Search API
Flight Search API is a Go + Beego service backed by Elasticsearch. It provides one search endpoint for querying flights with validation, sorting and paging support.
## Content Outline
1.  [Overview](#overview)
2.  [Features](#features)
3.  [Tech Stack](#tech-stack)
4.  [Prerequisites](#prerequisites)
5.  [Installation and Run](#installation-and-run)
6.  [Configuration](#configuration)
7.  [API Endpoint](#api-endpoint)
8.  [Search Parameters](#search-parameters)
9.  [Response and Error Format](#response-and-error-format)
10. [HTTP cURL Examples](#http-curl-examples)
11. [Project Folder Structure](#project-folder-structure)
## Overview
The API reads flight data from the Elasticsearch sample flights index and exposes a REST endpoint:
- Base URL: `/flight-api/v1`
- Search route: `GET /search`
## Features
- Flight search with multiple optional filters
- Sort by supported fields with order control
- Pagination support via limit and page parameters.
- Date, price, boolean, text and geo filter support.
## Tech Stack
- Go `1.25`
- Beego v2
- Elasticsearch v9 and Kibana v9 (via Docker Compose)

## Quick Setup 
### 1. Clone the project
```bash
git clone https://github.com/Rai321han/flight-api.git
cd flight-api
```
### 3. Configuration
```bash
mv conf/app.conf.example conf/app.conf
```
`conf/app.conf`:
```ini
appname = flight-api
httpport = 8080
runmode = dev
```
### 4. Start Elasticsearch and Kibana
```bash
docker compose up -d
```
This uses `docker-compose.yml` and starts:
- Elasticsearch on `http://localhost:9200`
- Kibana on `http://localhost:5601`
### 5. Install Go dependencies
```bash
go mod download
```
### 6. Run the API
```bash
go run main.go
```
Or 
```bash
bee run
```

API will start on:
- `http://localhost:8080`


## API Endpoint
### Search Flights
- Method: `GET`
- Path: `/flight-api/v1/search`
- Full local URL: `http://localhost:8080/flight-api/v1/search`
## Search Parameters
All parameters are optional.
| Parameter         | Type   | Description                          | Validation / Allowed Values                |
| ----------------- | ------ | ------------------------------------ | ------------------------------------------ |
| `carrier`         | string | Airline carrier                      | non-empty string                           |
| `flightNum`       | string | Flight number                        | non-empty string                           |
| `destCountry`     | string | Destination country code/name        | non-empty string                           |
| `originCountry`   | string | Origin country code/name             | non-empty string                           |
| `destCity`        | string | Destination city                     | non-empty string                           |
| `originCity`      | string | Origin city                          | non-empty string                           |
| `destAirportID`   | string | Destination airport ID               | non-empty string                           |
| `originAirportID` | string | Origin airport ID                    | non-empty string                           |
| `destAirport`     | string | Destination airport name             | non-empty string                           |
| `originAirport`   | string | Origin airport name                  | non-empty string                           |
| `dateFrom`        | date   | Start date filter                    | format `YYYY-MM-DD`                        |
| `dateTo`          | date   | End date filter                      | format `YYYY-MM-DD`, must be `>= dateFrom` |
| `priceMin`        | float  | Minimum ticket price                 | number, `>= 0`                             |
| `priceMax`        | float  | Maximum ticket price                 | number, `>= 0`, must be `>= priceMin`      |
| `cancelled`       | bool   | Cancelled flights filter             | `true` or `false`                          |
| `destLat`         | float  | Destination latitude for geo filter  | `-90` to `90`, requires `destLon`          |
| `destLon`         | float  | Destination longitude for geo filter | `-180` to `180`, requires `destLat`        |
| `originLat`       | float  | Origin latitude for geo filter       | `-90` to `90`, requires `originLon`        |
| `originLon`       | float  | Origin longitude for geo filter      | `-180` to `180`, requires `originLat`      |
| `sortBy`          | string | Sort field                           | `timestamp`, `AvgTicketPrice`              |
| `order`           | string | Sort direction                       | `asc`, `desc`                              |
| `limit`           | int    | Max records returned                 | positive integer, max `1000`, default `10` |
## Response and Error Format
### Success response
```json
{
  "data": {
    "flights": [
      {
        "Carrier": "Kibana Airlines",
        "FlightNum": "LY123"
      }
    ],
    "total": 13059,
    "limit": 1,
    "page": 1
  }
}
```
### Error response
```json
{
  "error": "error message"
}
```
## HTTP cURL Examples
```bash
# 1) Default search
curl -s "http://localhost:8080/flight-api/v1/search" | jq .
# 2) Filter by carrier
curl -s "http://localhost:8080/flight-api/v1/search?carrier=Kibana%20Airlines&limit=5" | jq .
# 3) Price range with descending sort
curl -s "http://localhost:8080/flight-api/v1/search?priceMin=100&priceMax=500&sortBy=AvgTicketPrice&order=desc&limit=5" | jq .
# 4) Cancelled flights only
curl -s "http://localhost:8080/flight-api/v1/search?cancelled=true&limit=5" | jq .
# 5) Date range
curl -s "http://localhost:8080/flight-api/v1/search?dateFrom=2024-01-01&dateTo=2024-03-31&limit=5" | jq .
# 6) Destination country filter
curl -s "http://localhost:8080/flight-api/v1/search?destCountry=CN&limit=5" | jq .
# 7) Validation error: inverted price range
curl -s "http://localhost:8080/flight-api/v1/search?priceMin=500&priceMax=100" | jq .
# 8) Validation error: bad date format
curl -s "http://localhost:8080/flight-api/v1/search?dateFrom=99-99-9999" | jq .
```
## Project Folder Structure
```text
.
├── conf/
│   └── app.conf                  # Beego app name, port, run mode
├── controllers/
│   └── search.go                 # HTTP handler for GET /flight-api/v1/search
├── models/
│   ├── app_error.go              # Typed AppError used across layers
│   ├── filters.go                # FlightFilters request model
│   ├── flight.go                 # Flight domain model(s)
│   └── service.go                # Service interface used by controller
├── routers/
│   └── router.go                 # Beego route and namespace registration
├── services/
│   ├── elasticsearch_client.go   # Elasticsearch client initialization
│   ├── flight_search_service.go  # Search business logic + ES call
│   ├── flight_service.go         # Service contract in services package
│   └── query_builder.go          # ES DSL query builder
├── utils/
│   └── response_format.go        # JSON response helpers
├── validators/
│   └── parameter.go              # Query parameter parsing and validation
├── docker-compose.yml            # Local Elasticsearch and Kibana stack
├── go.mod                        # Go module + dependencies
├── main.go                       # App entry point
├── .gitignore
└── README.md                     # Project documentation
```
