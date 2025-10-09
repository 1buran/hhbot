package apiclient

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/apiendpoint"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
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
func (ac ApiClient) SearchVacancies(text string, salary int) ([]dto.Vacancy, error) {
	var data dto.ResponseItems[dto.Vacancy]
	vacanciesEndpoint := apiendpoint.NewGetVacancies(text, salary)
	return ExtractItems(ac, vacanciesEndpoint, data)
}

// Get dictionary.
func (ac ApiClient) GetDictionary() (dict dto.Dictionary, err error) {
	dictionaryendpoint := apiendpoint.NewGetDictionary()
	err = Get(ac, dictionaryendpoint, &dict)
	return
}

// Apply to vacancy.
func (ac ApiClient) Apply(vacancyID, resumeID, message string) (*dto.Response, error) {
	applyendpoint := apiendpoint.NewPostNegotinations(vacancyID, resumeID, message)
	return Post(ac, applyendpoint)
}

// Get vacancy.
func (ac ApiClient) GetVacancy(vacancyID string) (vac dto.Vacancy, err error) {
	getvacancyendpoint := apiendpoint.NewGetVacancy(vacancyID)
	err = Get(ac, getvacancyendpoint, &vac)
	return
}

// New Api Client.
func NewApiClient(accessToken string) *ApiClient {
	return &ApiClient{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		accessToken: accessToken,
		userAgent:   "HHBot/1.0",
	}
}
