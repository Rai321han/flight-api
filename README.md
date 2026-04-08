# Flight Search API with Elasticsearch


## Overview

This project is a **flight search API** built with **Beego** and **Elasticsearch**, using the `kibana_sample_flights_data` dataset.  
It supports **advanced filtering, sorting, pagination, and aggregations**, enabling flexible and fast queries on flight data.

Both developers contribute **equally to Elasticsearch and API logic**, ensuring hands-on experience with queries, filters, and aggregations.

---

## Features

- Connect to Elasticsearch `kibana_sample_flights_data` index
- Search flights by:
  - Origin, Destination
  - Carrier
  - Cancellation status
  - Flight delay ranges
  - and more
- Aggregations:
  - Flights per Origin/Destination
  - Average ticket price ranges
- Pagination and sorting
- Full API endpoints with request validation
- Unit and integration tests with GoConvey & GoMonkey

---

## Tech Stack

- **Backend Framework:** Beego (Go)
- **Database:** Elasticsearch (sample dataset `kibana_sample_flights_data`)
- **Testing:** GoConvey, GoMonkey
- **Tools:** Dockers

---

## Test Search Queries

```bash 
# Default
curl -s "http://localhost:8080/flight-api/v1/search" | jq .
 
# By carrier
curl -s "http://localhost:8080/flight-api/v1/search?carrier=Kibana%20Airlines&limit=5" | jq .
 
# Price range sorted desc
curl -s "http://localhost:8080/flight-api/v1/search?priceMin=100&priceMax=500&sortBy=AvgTicketPrice&order=desc&limit=5%22 | jq .
 
# Cancelled only
curl -s "http://localhost:8080/flight-api/v1/search?cancelled=true&limit=5" | jq .
 
# Date range
curl -s "http://localhost:8080/flight-api/v1/search?dateFrom=2024-01-01&dateTo=2024-03-31&limit=5" | jq .
 
# Destination country
curl -s "http://localhost:8080/flight-api/v1/search?destCountry=CN&limit=5" | jq .
 
# Validation — inverted price
curl -s "http://localhost:8080/flight-api/v1/search?priceMin=500&priceMax=100" | jq .
 
# Validation — bad date
curl -s "http://localhost:8080/flight-api/v1/search?dateFrom=99-99-9999" | jq .
```
