package apiclient

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/1buran/hhbot/internal/infrastructure/apiclient/apiendpoint"
)

const (
	BaseURL = "https://api.hh.ru"
)

type ApiClient struct {
	httpClient  *http.Client
	accessToken string
	userAgent   string
}

func (ac ApiClient) CreateRequest(method, endpointUrl string, body io.Reader) (*http.Request, error) {
	u, err := url.JoinPath(BaseURL, endpointUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ac.userAgent)
	if ac.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+ac.accessToken)
	}
	return req, nil
}

// Search for vacancies.
func (ac ApiClient) SearchVacancies(text string, salary int) ([]VacancyDTO, error) {
	var data ResponseItemsDTO[VacancyDTO]
	vacanciesEndpoint := apiendpoint.NewGetVacancies(text, salary)
	return ExtractItems(ac, vacanciesEndpoint, data)
}

// New Api Client.
func NewApiClient(accessToken string) *ApiClient {
	return &ApiClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		accessToken: accessToken,
		userAgent:   "HHBot/1.0",
	}
}
