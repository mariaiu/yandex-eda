# Yandex.Eda Parsing Service

The Yandex.Eda Parsing Service is a versatile tool designed to parse restaurant data from Yandex.Food, allowing users to configure its behavior through a simple .yaml configuration file. The service offers powerful features for extracting restaurant information, managing ratings, and handling menu items.

## Table of Contents

- [Introduction](#yandexeda-parsing-service)
- [Features](#features)
- [Configuration](#configuration)
- [Usage](#usage)
  - [HTTP Requests](#http-requests)
  - [gRPC Calls](#grpc-calls)
- [Deployment](#deployment)

## Features

- **Configurability**: The service's behavior can be tailored using a .yaml configuration file. Users can specify:
  - Minimum restaurant ratings
  - Coordinates
  - API endpoint

- **Error Handling**: Invalid configuration files prevent the service from launching, ensuring proper functionality and user-friendly error messages.

- **Dynamic Rating Adjustment**: The minimum restaurant rating can be modified during application startup via configuration file flags.

- **Priority Coordinates**: High-priority coordinates can be specified using flags, overriding those in the configuration file.

- **Data Retrieval**: Upon a GET request to `/restaurant`, the service returns all restaurants from the database sorted by rating and price-to-weight ratio.

- **Detailed Menu Information**: A GET request to `/restaurant/:id` provides comprehensive pricing details for Philadelphia rolls with salmon at the selected restaurant.

- **Parallel Parsing**: With a GET request to `/parse` and basic authorization, the service performs multi-threaded parsing of restaurants meeting or exceeding the specified rating.

- **Database Interaction**: Parsed restaurants with relevant menu items are stored in the database for future reference.

- **Concurrency Control**: Prevents simultaneous parsing operations and enforces a timeout if parsing locks are not released within 30 seconds.

- **Optional Coordinates**: Coordinates in requests are optional. If absent, the service utilizes coordinates provided in application launch flags.

- **Logging**: Parsing activities are logged to the standard output using the `github.com/sirupsen/logrus v1.6.0` logger.

## Configuration

The service's behavior can be configured through the `.yaml` file. Configure the following parameters:

- Minimum restaurant ratings
- Coordinates
- API endpoint

## Usage

### HTTP Requests

Make HTTP requests to interact with the service. Examples include:

```bash
# Parse restaurants with specified workers and coordinates
curl -u alice:alice 'localhost:8080/parse?workers=25&latitude=59.836685&longitude=30.358017'

# Retrieve detailed menu information for a specific restaurant
curl localhost:8080/restaurant/1 | jq

# Retrieve a list of restaurants sorted by rating and price-to-weight ratio
curl localhost:8080/restaurant | jq

```
### gRPC Calls
For gRPC calls, use the following commands:

```bash
# Start a gRPC session
evans proto/restaurant.proto -p 8081

# Call the ParseRestaurants method
call ParseRestaurants

# Call the GetRestaurant method
call GetRestaurant

# Call the GetRestaurants method
call GetRestaurants

```

## Deployment

- Run migrations: make migrate
- Generate proto files: make proto-gen
- Start the application: make docker-up
