package country

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/itsdarkhost/rbk-week2/internal/utils"
)

var ErrCountryNotFound = errors.New("country not found")

type Client struct {
	api        *utils.Client
	metadata   *utils.Client
	countries  []Country
	codeByName map[string]string
}

func NewClient(httpClient *http.Client, baseURL, metadataBaseURL string) *Client {
	return &Client{
		api:        utils.NewClient(baseURL, httpClient, nil),
		metadata:   utils.NewClient(metadataBaseURL, httpClient, nil),
		codeByName: map[string]string{},
	}
}

func (c *Client) FindCountry(ctx context.Context, country string) (Country, error) {
	countries, err := c.listCountries(ctx)
	if err != nil {
		return Country{}, err
	}

	target := normalize(country)
	var fallback Country
	var hasFallback bool

	for _, item := range countries {
		current := normalize(item.Name)
		if current == target {
			return item, nil
		}

		if !hasFallback && (strings.Contains(current, target) || strings.Contains(target, current)) {
			fallback = item
			hasFallback = true
		}
	}

	if hasFallback {
		return fallback, nil
	}

	return Country{}, ErrCountryNotFound
}

func (c *Client) CountryCode(ctx context.Context, country string) (string, error) {
	normalized := normalize(country)

	if code, ok := c.codeByName[normalized]; ok {
		return code, nil
	}

	code, err := c.fetchCountryCode(ctx, country, true)
	if err != nil {
		code, err = c.fetchCountryCode(ctx, country, false)
		if err != nil {
			return "", err
		}
	}

	c.codeByName[normalized] = code

	return code, nil
}

func (c *Client) listCountries(ctx context.Context) ([]Country, error) {
	if len(c.countries) > 0 {
		cached := append([]Country(nil), c.countries...)
		return cached, nil
	}

	resp, err := c.api.Get(ctx, "/countries", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, utils.ReadErrorResponse(resp)
	}

	var payload countriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	countries := make([]Country, 0, len(payload.Data))
	for _, item := range payload.Data {
		countries = append(countries, Country{
			Name:   item.Country,
			Cities: item.Cities,
		})
	}

	c.countries = countries

	return append([]Country(nil), countries...), nil
}

func (c *Client) fetchCountryCode(ctx context.Context, country string, fullText bool) (string, error) {
	query := map[string]string{
		"fields": "name,cca2",
	}
	if fullText {
		query["fullText"] = "true"
	}

	resp, err := c.metadata.Get(ctx, "/v3.1/name/"+url.PathEscape(country), query)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrCountryNotFound
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return "", utils.ReadErrorResponse(resp)
	}

	var payload metadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if len(payload) == 0 {
		return "", ErrCountryNotFound
	}

	target := normalize(country)
	for _, item := range payload {
		if normalize(item.Name.Common) == target || normalize(item.Name.Official) == target {
			return item.CCA2, nil
		}
	}

	return payload[0].CCA2, nil
}

func normalize(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer("-", " ", "_", " ", "'", "", ".", "", ",", "").Replace(value)
	return strings.Join(strings.Fields(value), " ")
}
