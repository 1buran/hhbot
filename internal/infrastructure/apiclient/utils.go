package apiclient

import (
	"encoding/json"
	"fmt"

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

// Generic GET request.
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

type ErrorDetails struct {
	Type  string
	Value string
}

type BadRequestData struct {
	Request_id  string
	Description string
	Errors      []ErrorDetails
}

// Generic POST request.
func Post(
	ac ApiClient,
	endpoint apiendpoint.ApiEndpoint,
) (*dto.Response, error) {
	req, err := ac.CreateRequest("POST", endpoint.Url(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = endpoint.Payload().Encode()

	res, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var msg BadRequestData
	if res.StatusCode >= 400 && res.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(res.Body).Decode(&msg); err != nil {
			return nil, err
		}
	}

	var errMsg string
	if msg.Description != "" {
		errMsg = fmt.Sprintf("Description: %s", msg.Description)
	}
	for _, e := range msg.Errors {
		errMsg += fmt.Sprintf("\n%s: %s", e.Type, e.Value)
	}

	return &dto.Response{
		Code:    res.StatusCode,
		Status:  res.Status,
		Headers: res.Header,
		Error:   errMsg,
	}, nil
}
