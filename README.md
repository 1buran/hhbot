# HH.ru API Client

A Go client library for the HeadHunter (hh.ru) API, generated from OpenAPI specification using oapi-codegen.

## Features

This client provides access to key HeadHunter API operations:

- **Search Vacancies** - Find job vacancies by various criteria (location, salary, keywords, etc.)
- **Get Vacancy Details** - Retrieve detailed information about specific vacancies
- **Apply to Vacancy** - Submit job applications to vacancies

## Prerequisites

- Go 1.25.1 or later
- Internet connection to access hh.ru API

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd hhbot
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate the client code:
   ```bash
   go generate ./...
   ```

## Configuration

The client generation is configured in `cfg.yaml`:

```yaml
package: client
output-options:
  include-operation-ids:
    - get-vacancies      # Search for vacancies
    - get-vacancy        # Get vacancy details
    - apply-to-vacancy   # Apply to vacancy
generate:
  models: true
  client: true
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/1buran/hhbot/client"
)

func main() {
    // Create HTTP client
    httpClient := &http.Client{}
    
    // Create API client
    apiClient, err := client.NewClient("https://api.hh.ru", client.WithHTTPClient(httpClient))
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    
    ctx := context.Background()
    // Use the client...
}
```

### Search for Vacancies

```go
// Search parameters
params := &client.GetVacanciesParams{
    Text:           stringPtr("golang developer"),  // Search query
    Area:           stringPtr("1"),                 // Moscow (area ID)
    Page:           float32Ptr(0),                  // Page number
    PerPage:        float32Ptr(20),                 // Results per page
    Salary:         float32Ptr(100000),             // Minimum salary
    Currency:       stringPtr("RUR"),               // Currency
    OnlyWithSalary: boolPtr(true),                  // Only with salary
}

response, err := apiClient.GetVacancies(ctx, params)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
defer response.Body.Close()
```

### Get Vacancy Details

```go
vacancyID := "123456"
response, err := apiClient.GetVacancy(ctx, vacancyID, &client.GetVacancyParams{})
if err != nil {
    log.Printf("Error: %v", err)
    return
}
defer response.Body.Close()
```

### Apply to Vacancy

```go
body := client.ApplyToVacancyMultipartRequestBody{
    // Add required fields based on your needs
    // VacancyId: "123456",
    // ResumeId: "your_resume_id",
    // Message: "Your cover letter",
}

response, err := apiClient.ApplyToVacancy(ctx, body)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
defer response.Body.Close()
```

## API Reference

### Search Parameters

The `GetVacanciesParams` struct supports extensive filtering options:

- `Text` - Search query text
- `Area` - Region/city ID (use `/areas` endpoint to get IDs)
- `Professional_role` - Professional category ID
- `Industry` - Industry ID
- `Employer_id` - Specific employer ID
- `Salary` - Minimum salary amount
- `Currency` - Currency code (RUR, USD, EUR, etc.)
- `Experience` - Required experience level
- `Employment` - Employment type
- `Schedule` - Work schedule
- `OnlyWithSalary` - Show only vacancies with salary specified
- `Period` - Search within last N days
- `DateFrom`/`DateTo` - Date range filters
- Geographic coordinates for location-based search
- Sorting and pagination options

### Response Handling

All API methods return HTTP responses. Check the status code and parse the response body according to your needs:

```go
switch response.StatusCode {
case 200:
    // Success - parse response body
case 400:
    // Bad request - check parameters
case 403:
    // Authorization required or forbidden
case 404:
    // Not found
case 429:
    // Rate limit exceeded
default:
    // Handle other status codes
}
```

## Authentication

For full API access, you'll need to implement OAuth2 authentication with HeadHunter. Add the authorization token to your HTTP client:

```go
httpClient := &http.Client{
    Transport: &authTransport{
        token: "your_access_token",
        base:  http.DefaultTransport,
    },
}
```

## Rate Limiting

The HeadHunter API has rate limits. Implement appropriate retry logic and respect the rate limits to avoid being blocked.

## Examples

See `example_usage.go` for complete working examples of all supported operations.

## Development

### Regenerating the Client

If you modify `cfg.yaml` or update the OpenAPI specification:

```bash
go generate ./...
```

### Adding New Operations

1. Find the operation ID in `openapi.yml`
2. Add it to the `include-operation-ids` list in `cfg.yaml`
3. Regenerate the client

## License

See LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For API documentation and support, visit:
- [HeadHunter API Documentation](https://github.com/hhru/api)
- [HeadHunter Developer Portal](https://dev.hh.ru/)
