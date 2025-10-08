# HH.ru Vacancy Search Bot with OAuth

A minimal Go client for searching vacancies on HeadHunter (hh.ru) with OAuth 2.0 authentication.

## Features

- OAuth 2.0 user authentication with CSRF protection (state parameter)
- Search vacancies by keywords, salary, location, employment type
- Filter by work schedule and experience
- Automatic token exchange and management
- Local HTTP server for OAuth callback handling

## Setup

### 1. Register Your Application

1. Go to [https://dev.hh.ru/](https://dev.hh.ru/)
2. Click "Создать приложение" (Create Application)
3. Fill in the application details:
   - **Name**: Your bot name (e.g., "Vacancy Search Bot")
   - **Redirect URI**: `http://localhost:8080/callback`
   - **Description**: Brief description of your application
4. Save your **Client ID** and **Client Secret**

### 2. Create .env File

Create a `.env` file in the project root:

```bash
HH_CLIENT_ID=your_client_id_here
HH_CLIENT_SECRET=your_client_secret_here
```

You can copy from the example:
```bash
cp .env.example .env
# Then edit .env with your credentials
```

### 3. Install and Run

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go
```

## How It Works

### OAuth Flow

1. Application loads credentials from `.env` file
2. Generates random `state` parameter for CSRF protection
3. Starts local HTTP server on `http://localhost:8080`
4. Constructs authorization URL:
   ```
   https://hh.ru/oauth/authorize?
   response_type=code&
   client_id={client_id}&
   state={state}&
   redirect_uri=http://localhost:8080/callback
   ```
5. User visits URL and authorizes the application
6. HH.ru redirects to callback with `code` and `state` parameters
7. Application verifies `state` matches (CSRF protection)
8. Exchanges authorization code for access token
9. Token is used for API requests

### Search Parameters

The client supports the following search parameters:

- **Text** - Keywords to search in vacancy (e.g., "golang developer")
- **Area** - Location ID (1 = Moscow, 2 = St. Petersburg)
- **Salary** - Minimum salary amount
- **Currency** - Currency code (RUR, USD, EUR)
- **Employment** - Employment type:
  - `full` - Full-time
  - `part` - Part-time
  - `project` - Project work
  - `volunteer` - Volunteer
  - `probation` - Probation
- **Schedule** - Work schedule:
  - `fullDay` - Full day
  - `shift` - Shift work
  - `flexible` - Flexible schedule
  - `remote` - Remote work
  - `flyInFlyOut` - Fly-in fly-out
- **OnlyWithSalary** - Show only vacancies with specified salary
- **Page** - Page number (0-based)
- **PerPage** - Results per page (max 100)

## Usage Examples

### Basic Search

```go
params := hhclient.SearchParams{
    Text:    "golang developer",
    Area:    "1", // Moscow
    PerPage: 20,
}

result, err := client.SearchVacancies(params)
```

### Search with Salary Filter

```go
params := hhclient.SearchParams{
    Text:           "python developer",
    Salary:         150000,
    Currency:       "RUR",
    OnlyWithSalary: true,
}
```

### Search Remote Jobs

```go
params := hhclient.SearchParams{
    Text:       "backend developer",
    Schedule:   "remote",
    Employment: "full",
}
```

### Search by Multiple Criteria

```go
params := hhclient.SearchParams{
    Text:           "DevOps engineer",
    Area:           "1",
    Salary:         200000,
    Currency:       "RUR",
    Employment:     "full",
    Schedule:       "remote",
    OnlyWithSalary: true,
    PerPage:        50,
}
```

## Security Features

### CSRF Protection

The application implements CSRF protection using the `state` parameter:

1. Generates a random state value before authorization
2. Includes state in authorization URL
3. HH.ru returns the same state in callback
4. Application verifies state matches before exchanging code
5. Rejects requests with mismatched state values

This prevents cross-site request forgery attacks during the OAuth flow.

## API Reference

### Client Methods

#### `NewClient(accessToken string) *Client`
Creates a new API client with the given access token.

#### `SearchVacancies(params SearchParams) (*SearchResponse, error)`
Searches for vacancies based on the provided parameters.

### OAuth Methods

#### `GenerateState() (string, error)`
Generates a random state string for CSRF protection.

#### `GetAuthURL(state string) string`
Returns the authorization URL with state parameter.

#### `ExchangeCode(code, state string) (*TokenResponse, error)`
Exchanges the authorization code for an access token.

## Response Structure

```go
type SearchResponse struct {
    Items   []Vacancy  // List of vacancies
    Found   int        // Total number of found vacancies
    Pages   int        // Total number of pages
    Page    int        // Current page number
    PerPage int        // Items per page
}

type Vacancy struct {
    ID         string
    Name       string
    Area       Area
    Salary     *Salary
    Employer   Employer
    Snippet    Snippet
    Employment *IdName
    Schedule   *IdName
    URL        string
}
```

## Common Area IDs

- `1` - Moscow
- `2` - St. Petersburg
- `113` - Russia (all regions)
- `1001` - Other countries

For complete list, use the `/areas` API endpoint.

## Troubleshooting

### "Please set HH_CLIENT_ID and HH_CLIENT_SECRET in .env file"
Create a `.env` file in the project root with your credentials.

### "Authorization failed"
- Check that redirect URI matches exactly: `http://localhost:8080/callback`
- Ensure port 8080 is not in use by another application
- Verify Client ID and Secret are correct in `.env` file

### "state mismatch: possible CSRF attack"
This indicates a security issue. Make sure:
- You're using the exact URL provided by the application
- No other application is interfering with the OAuth flow
- Try restarting the application to generate a new state

### "API error 403"
Your access token may have expired. Re-run the application to get a new token.

### "API error 400"
Check your search parameters. Some combinations may be invalid.

## License

See LICENSE file for details.
