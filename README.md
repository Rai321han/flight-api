# Flight Search API with Elasticsearch

testing git&github

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
