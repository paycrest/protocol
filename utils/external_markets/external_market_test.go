package externalmarkets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExternalMarketRates(t *testing.T) {
	// Mock Bitget server with currency-specific responses
	bitgetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("fiat")
		var response string

		switch currency {
		case "NGN":
			response = `{
                "code": "00000",
                "msg": "success",
                "data": [
                    {"price": "740.0", "available": "1000"},
                    {"price": "745.0", "available": "2000"},
                    {"price": "750.0", "available": "1500"}
                ]
            }`
		case "KES":
			response = `{
                "code": "00000",
                "msg": "success",
                "data": [
                    {"price": "145.0", "available": "1000"},
                    {"price": "145.5", "available": "2000"},
                    {"price": "146.0", "available": "1500"}
                ]
            }`
		case "GHS":
			response = `{
                "code": "00000",
                "msg": "success",
                "data": [
                    {"price": "545.0", "available": "1000"},
                    {"price": "545.5", "available": "2000"},
                    {"price": "546.0", "available": "1500"}
                ]
            }`
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer bitgetServer.Close()

	// Mock Binance server with currency-specific responses
	binanceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		var response string

		switch symbol {
		case "USDTKES":
			response = `{"symbol":"USDTKES","price":"145.50"}`
		case "USDTGHS":
			response = `{"symbol":"USDTGHS","price":"545.50"}`
		default:
			response = `{"symbol":"USDTNGN","price":"750.00"}`
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer binanceServer.Close()

	// Mock Quidax server with currency-specific responses
	quidaxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
            "data": {
                "last_price": "755.00"
            }
        }`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer quidaxServer.Close()

	// Create ExternalMarketRates instance with mock servers
	emr := NewExternalMarketRates()
	emr.bitgetURL = bitgetServer.URL
	emr.binanceURL = binanceServer.URL
	emr.quidaxURL = quidaxServer.URL

	ctx := context.Background()

	t.Run("FetchRate", func(t *testing.T) {
		tests := []struct {
			name     string
			currency string
			want     float64
			wantErr  bool
		}{
			{
				name:     "Test NGN rate fetch",
				currency: "NGN",
				want:     750.0,
				wantErr:  false,
			},
			{
				name:     "Test other currency rate fetch (KES)",
				currency: "KES",
				want:     145.5,
				wantErr:  false,
			},
			{
				name:     "Test other currency rate fetch (GHS)",
				currency: "GHS",
				want:     545.5,
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := emr.FetchRate(ctx, tt.currency)
				if (err != nil) != tt.wantErr {
					t.Errorf("FetchRate() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.Price != tt.want {
					t.Errorf("FetchRate() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("CalculateMedian", func(t *testing.T) {
		tests := []struct {
			name   string
			values []float64
			want   float64
		}{
			{
				name:   "Odd number of values",
				values: []float64{1, 2, 3},
				want:   2,
			},
			{
				name:   "Even number of values",
				values: []float64{1, 2, 3, 4},
				want:   2.5,
			},
			{
				name:   "Empty slice",
				values: []float64{},
				want:   0,
			},
			{
				name:   "Single value",
				values: []float64{1},
				want:   1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := calculateMedian(tt.values); got != tt.want {
					t.Errorf("calculateMedian() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
