package utils

import (
	"fmt"

	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/paycrest/aggregator/config"
)

// Provider represents a rate provider
type Provider string

const (
	ProviderBitget  Provider = "BITGET"
	ProviderBinance Provider = "BINANCE"
	ProviderQuidax  Provider = "QUIDAX"
)

// Rate holds currency pair information *****
type Rate struct {
	Currency  string
	Price     float64
	Provider  Provider
	Timestamp time.Time
}

// RateResponse represents a standardized rate response
type RateResponse struct {
	Rate  Rate
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
	quidaxURL  string
	bitgetURL  string
	binanceURL string
}

// NewExternalMarketRates creates a new instance of ExternalMarketRates
func NewExternalMarketRates() *ExternalMarketRates {
	cfg := config.ServerConfig()
	return &ExternalMarketRates{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		quidaxURL:  cfg.QuidaxURL,
		bitgetURL:  cfg.BitgetURL,
		binanceURL: cfg.BinanceURL,
	}
}

// FetchRate gets the median rate for a given currency
func (e *ExternalMarketRates) FetchRate(ctx context.Context, currency string) (Rate, error) {
	switch strings.ToUpper(currency) {
	case "NGN":
		return e.fetchNGNRate(ctx)
	default:
		return e.fetchOtherCurrencyRate(ctx, currency)
	}
}

// fetchNGNRate fetches rate for NGN from Quidax and BItget
func (e *ExternalMarketRates) fetchNGNRate(ctx context.Context) (Rate, error) {
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
	var validRates []Rate
	for i := 0; i < 2; i++ {
		resp := <-rates
		if resp.Error == nil {
			validRates = append(validRates, resp.Rate)
		}
	}
	if len(validRates) == 0 {
		return Rate{}, fmt.Errorf("no valid rates found for NGN")
	}
	return calculateMedianRate(validRates), nil
}

// fetchOtherCurrencyRate fetches rate for other currencies from Binance and Bidget
func (e *ExternalMarketRates) fetchOtherCurrencyRate(ctx context.Context, currency string) (Rate, error) {
	rates := make(chan RateResponse, 2)

	go func() {
		rate, err := e.fetchBinanceRate(ctx, currency)
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	go func() {
		rate, err := e.fetchBitgetRate(ctx, currency)
		rates <- RateResponse{Rate: rate, Error: err}
	}()

	var validRates []Rate
	for i := 0; i < 2; i++ {
		resp := <-rates
		if resp.Error == nil {
			validRates = append(validRates, resp.Rate)
		}
	}
	if len(validRates) == 0 {
		return Rate{}, fmt.Errorf("no valid rates found for %s", currency)
	}
	return calculateMedianRate(validRates), nil
}

// fetchBitgetRate fetches rates from Bitget
func (e *ExternalMarketRates) fetchBitgetRate(ctx context.Context, currency string) (Rate, error) {
	url := fmt.Sprintf("%s?symbol=USDT%s&limit=20", e.bitgetURL, strings.ToUpper(currency))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Rate{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return Rate{}, fmt.Errorf("fetching from Bitget:: %w", err)
	}
	defer resp.Body.Close()

	var result BitgetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Rate{}, fmt.Errorf("decoding Bitget response: %w", err)
	}

	if len(result.Data) == 0 {
		return Rate{}, fmt.Errorf("no Bitget P2P ads found for %s", currency)
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
		return Rate{}, fmt.Errorf("no valid Bitget P2P prices found for %s", currency)
	}

	medianPrice := calculateMedian(prices)
	return Rate{
		Currency:  currency,
		Price:     medianPrice,
		Provider:  ProviderBitget,
		Timestamp: time.Now(),
	}, nil

}

// fetchBinanceRate fetches rates from Binance
func (e *ExternalMarketRates) fetchBinanceRate(ctx context.Context, currency string) (Rate, error) {
	url := fmt.Sprintf("%s?symbol=USDT%s", e.binanceURL, strings.ToUpper(currency))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Rate{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return Rate{}, fmt.Errorf("fetching from Binance:: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Rate{}, fmt.Errorf("decoding Binance response: %w", err)
	}
	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return Rate{}, fmt.Errorf("parsing Binance price: %w", err)
	}
	return Rate{
		Currency:  currency,
		Price:     price,
		Provider:  ProviderBinance,
		Timestamp: time.Now(),
	}, nil
}

// fetchQuidaxRate fetches ratesfrom Quidax
func (e *ExternalMarketRates) fetchQuidaxRate(ctx context.Context, currency string) (Rate, error) {
	url := fmt.Sprintf("%s/usdt%s/ticker", e.quidaxURL, strings.ToLower(currency))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Rate{}, fmt.Errorf("creating request: %w", err)
	}
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return Rate{}, fmt.Errorf("fetching from Quidax:: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			LastPrice string `json:"last_price"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Rate{}, fmt.Errorf("decoding Quidax response: %w", err)
	}

	price, err := strconv.ParseFloat(result.Data.LastPrice, 64)
	if err != nil {
		return Rate{}, fmt.Errorf("parsing Quidax price: %w", err)
	}

	return Rate{
		Currency:  currency,
		Price:     price,
		Provider:  ProviderQuidax,
		Timestamp: time.Now(),
	}, nil
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

// New helper function to calculate median from Rate objects
func calculateMedianRate(rates []Rate) Rate {
	if len(rates) == 0 {
		return Rate{}
	}

	// Extract prices for median calculation
	prices := make([]float64, len(rates))
	for i, rate := range rates {
		prices[i] = rate.Price
	}

	medianPrice := calculateMedian(prices)

	// Find the rate object closest to the median price
	var closestRate Rate
	smallestDiff := float64(^uint(0) >> 1) // Max float64

	for _, rate := range rates {

		// Find closest rate to median
		diff := Abs(rate.Price - medianPrice)
		if diff < smallestDiff {
			smallestDiff = diff
			closestRate = rate
		}
	}

	// Return a new Rate with the median price but keeping other metadata from the closest rate
	return Rate{
		Currency:  closestRate.Currency,
		Price:     medianPrice,
		Provider:  closestRate.Provider,
		Timestamp: closestRate.Timestamp,
	}
}
