package apiclient

import (
	"crypto/rand"
	"encoding/base64"

	"encoding/json"

	"github.com/1buran/hhbot/internal/infrastructure/apiclient/apiendpoint"
)

// Generic extraction of items from response with items included.
func ExtractItems[T ResponseType](
	ac ApiClient,
	endpoint apiendpoint.ApiEndpoint,
	data ResponseItemsDTO[T],
) ([]T, error) {
	req, err := ac.CreateRequest("GET", endpoint.Url(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = endpoint.Payload().Encode()
	res, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Items, nil
}

// GenerateState creates a random state string for CSRF protection
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
