package externalmarkets

import (
	"fmt"

	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Provider represents a rate provider
type Provider string

// const (
// 	ProviderBitget  Provider = "BITGET"
// 	ProviderBinance Provider = "BINANCE"
// 	ProviderQuidax  Provider = "QUIDAX"
// )

// // Rate holds currency pair information *****
// type Rate struct {
// 	Currency  string
// 	Price     float64
// 	Provider  Provider
// 	Timestamp time.Time
// }

// RateResponse represents a standardized rate response
type RateResponse struct {
	Rate  float64
	Error error
}

// BitgetP2PAd represents a P2P advertisement from Bitget
type BitgetP2PAd struct {
	Price     string `json:"price"`
	Available string `json:"available"`
}

// BitgetResponse represents the API response from Bitget
type BitgetResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"msg"`
	Data    []BitgetP2PAd `json:"data"`
}

// ExternalMarketRates handles fetching and calculating rates from external providers
type ExternalMarketRates struct {
	httpClient *http.Client
	bitgetURL  string
	binanceURL string
	quidaxURL  string
}

// NewExternalMarketRates creates a new instance of ExternalMarketRates
func NewExternalMarketRates() *ExternalMarketRates {
	return &ExternalMarketRates{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		bitgetURL:  "https://api.bitget.com/api/v1/p2p/ads",
		binanceURL: "https://api.binance.com/api/v3/ticker/price",
		quidaxURL:  "https://www.quidax.com/api/v1/markets",
	}
}

// FetchRate gets the median rate for a given currency
func (e *ExternalMarketRates) FetchRate(ctx context.Context, currency string) (float64, error) {
	switch strings.ToUpper(currency) {
	case "NGN":
		return e.fetchNGNRate(ctx)
	default:
		return e.fetchOtherCurrencyRate(ctx, currency)
	}
}

// fetchNGNRate fetches rate for NGN from Quidax and BItget
func (e *ExternalMarketRates) fetchNGNRate(ctx context.Context) (float64, error) {
	rates := make(chan RateResponse, 2)

	// Fetch rates concurrently
	go func() {
		rate, err := e.fetchQuidaxRate(ctx, "NGN")
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	go func() {
		rate, err := e.fetchBitgetRate(ctx, "NGN")
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	// Collect results
	var validRates []float64
	for i := 0; i < 2; i++ {
		resp := <-rates
		if resp.Error == nil {
			validRates = append(validRates, resp.Rate)
		}
	}
	if len(validRates) == 0 {
		return 0, fmt.Errorf("no valid rates found for NGN")
	}
	return calculateMedian(validRates), nil
}

// fetchOtherCurrencyRate fetches rate for other currencies from Binance and Bidget
func (e *ExternalMarketRates) fetchOtherCurrencyRate(ctx context.Context, currency string) (float64, error) {
	rates := make(chan RateResponse, 2)

	go func() {
		rate, err := e.fetchBinanceRate(ctx, currency)
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	go func() {
		rate, err := e.fetchBitgetRate(ctx, currency)
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	var validRates []float64
	for i := 0; i < 2; i++ {
		resp := <-rates
		if resp.Error == nil {
			validRates = append(validRates, resp.Rate)
		}
	}
	if len(validRates) == 0 {
		return 0, fmt.Errorf("no valid rates found for %s", currency)
	}
	return calculateMedian(validRates), nil
}

// fetchBitgetRate fetches rates from Bitget
func (e *ExternalMarketRates) fetchBitgetRate(ctx context.Context, currency string) (float64, error) {
	url := fmt.Sprintf("%s?fiat=%s&crypto=USDT&type=BUY&limit=20", e.bitgetURL, currency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetching from Bitget:: %w", err)
	}
	defer resp.Body.Close()

	var result BitgetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decoding Bitget response: %w", err)
	}

	if len(result.Data) == 0 {
		return 0, fmt.Errorf("no Bitget P2P ads found for %s", currency)
	}

	// Extract and convert prices to float64
	var prices []float64
	for _, ad := range result.Data {
		price, err := strconv.ParseFloat(ad.Price, 64)
		if err != nil {
			continue
		}
		prices = append(prices, price)
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("no valid Bitget P2P prices found for %s", currency)
	}

	return calculateMedian(prices), nil

}

// fetchBinanceRate fetches rates from Binance
func (e *ExternalMarketRates) fetchBinanceRate(ctx context.Context, currency string) (float64, error) {
	symbol := fmt.Sprintf("USDT%s", currency)
	url := fmt.Sprintf("%s?symbol=%s", e.binanceURL, symbol)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetching from Binance:: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decoding Binance response: %w", err)
	}
	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing Binance price: %w", err)
	}
	return price, nil
}

// fetchQuidaxRate fetches ratesfrom Quidax
func (e *ExternalMarketRates) fetchQuidaxRate(ctx context.Context, currency string) (float64, error) {
	url := fmt.Sprintf("%s/usd%s/ticker", e.quidaxURL, strings.ToLower(currency))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetching from Quidax:: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			LastPrice string `json:"last_price"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decoding Quidax response: %w", err)
	}

	price, err := strconv.ParseFloat(result.Data.LastPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing Quidax price: %w", err)
	}
	return price, nil
}

// calculateMedian calculates the median value from slice of float64
func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sort.Float64s(values)
	middle := len(values) / 2

	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}
	return values[middle]
}
