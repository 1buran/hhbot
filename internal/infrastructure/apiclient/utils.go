package apiclient

import (
	"encoding/json"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/apiendpoint"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

// Generic extraction of items from response with items included.
func ExtractItems[T dto.ResponseType](
	ac ApiClient,
	endpoint apiendpoint.ApiEndpoint,
	data dto.ResponseItems[T],
) ([]T, error) {
	var items []T
	var totalPages int

	for {
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

		items = append(items, data.Items...)
		totalPages = data.Pages

		if data.Page == totalPages-1 { // page numbers started from 0
			break
		} else {
			endpoint.SetPage(data.Page + 1)
		}
	}
	return items, nil
}

// Generic Get request.
func Get[T dto.ResponseType](
	ac ApiClient,
	endpoint apiendpoint.ApiEndpoint,
	data *T,
) error {
	req, err := ac.CreateRequest("GET", endpoint.Url(), nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = endpoint.Payload().Encode()
	res, err := ac.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	return nil
}
