package fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/KDarenskii/worker-service/internal/app/config/section"
	"github.com/KDarenskii/worker-service/internal/app/entity"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const requestTimeout = 10 * time.Second

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type fixerResponse struct {
	Success bool               `json:"success"`
	Base    string             `json:"base"`
	Date    string             `json:"date"`
	Rates   map[string]float64 `json:"rates"`
	Error   *fixerError        `json:"error,omitempty"`
}

type fixerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func NewClient(cfg section.ClientFixer) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   requestTimeout,
		},
		apiKey:  cfg.ApiKey,
		baseURL: cfg.BaseURL,
	}
}

// GetRates возвращает курсы валют относительно базовой валюты base.
func (c *Client) GetRates(ctx context.Context, base string) (map[string]float64, error) {
	args := url.Values{}

	args.Set("base", base)
	args.Set("access_key", c.apiKey)

	requestUrl := fmt.Sprintf("%s/latest?%s", c.baseURL, args.Encode())

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerUnavailable, err.Error())
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerUnavailable, err.Error())
	}

	defer func() {
		if response != nil {
			_ = response.Body.Close()
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, entity.ErrFixerUnavailable
	}

	var fixerResp fixerResponse

	if err := json.NewDecoder(response.Body).Decode(&fixerResp); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerInvalidResponse, err.Error())
	}

	if !fixerResp.Success {
		return nil, mapFixerError(fixerResp.Error)
	}

	log.Info().Ctx(ctx).Str("base", base).Int("count", len(fixerResp.Rates)).Msg("Получены курсы валют от Fixer")

	return fixerResp.Rates, nil
}

func mapFixerError(fixerErr *fixerError) error {
	if fixerErr == nil {
		return entity.ErrFixerInvalidResponse
	}

	switch fixerErr.Code {
	case 101:
		return entity.ErrFixerInvalidApiKey
	case 104, 105:
		return entity.ErrFixerRateLimitExceeded
	default:
		return fmt.Errorf(
			"%w: [%d] %s - %s", entity.ErrFixerInvalidResponse,
			fixerErr.Code, fixerErr.Type, fixerErr.Info,
		)
	}
}
