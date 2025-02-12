package services

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/aggregator/config"
	"github.com/paycrest/aggregator/ent"
	"github.com/paycrest/aggregator/types"
	"github.com/paycrest/aggregator/utils"
	"github.com/paycrest/aggregator/utils/logger"
	tokenUtils "github.com/paycrest/aggregator/utils/token"
)

var streamConf = config.StreamConfig()


// StreamManager is an interface for managing streams
type StreamManager interface {
    // AddressStreamManager is an interface for managing address streams
    CreateAddressStream(ctx context.Context, network *ent.Network, startRange int, endRange int) (string, error)
    UpdateAddressStream(ctx context.Context, streamId string, network *ent.Network, startRange int, endRange int) error
    DeleteAddressStream(ctx context.Context, streamId string) error
	PauseAddressStream(ctx context.Context, streamId string) (error)
	ActivateAddressStream(ctx context.Context, streamId string) (error)
    
    // EventStreamManager is an interface for managing event streams
    // CreateContractEventStream(ctx context.Context, params types.StreamCreationParams) (string, error)
    // DeleteContractEventStream(ctx context.Context, streamID string) error
    // UpdateContractEventStream(ctx context.Context, streamID string, params types.StreamCreationParams) error
	// GetContractEventStream(ctx context.Context, streamID string) (*types.StreamReturnPayload, error)

	GetAllStreams(ctx context.Context) ([]*types.StreamReturnPayload, error)
}

// QuickNodeStreamManager is a service for managing streams using QuickNode
type QuickNodeStreamManager struct {
	streamManager StreamManager
}


// This function will need to be called everytime a new receive address is generated
func (q *QuickNodeStreamManager) CreateAddressStream(ctx context.Context, network *ent.Network, startRange int, endRange int, filterConfig types.FilterConfig) (string, error) {
	if endRange <= startRange {
		return "", fmt.Errorf("endRange must be greater than startRange")
	}
	params := types.StreamCreationParams{
		Name: "PaycrestAddressStream",
		Network: network.Identifier,
		Dataset: "block_with_receipts",
		FilterFunction: utils.GetEncodedFilterFunction(filterConfig),
		Region: "usa_east",
		StartRange: startRange,
		EndRange: endRange,
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		IncludeStreamMetadata: "body",
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: map[string]interface{}{
			"url": "webhook", // todo, set up a webhook to receive the data and use HMAC to verify the data
			"compression": "none",
			"max_retry": 5,
			"retry_interval_sec": 1,
			"post_timeout_sec":   10,
		},
	}

	streamId, err := q.createQuickNodeStream(ctx, params, streamConf.QuickNodeAPIURL)
	if err != nil {
		return "", err
	}
    return streamId, nil
}

func (q *QuickNodeStreamManager) UpdateAddressStream(ctx context.Context, streamId string, network *ent.Network, startRange int, endRange int, filterConfig types.FilterConfig) error {
	if endRange <= startRange {
		return fmt.Errorf("UPDATE_STREAM: endRange must be greater than startRange")
	}
	if streamId == "" {
		return fmt.Errorf("UPDATE_STREAM: streamId is required")
	}
	params := types.StreamCreationParams{
		Name: "PaycrestAddressStream",
		FilterFunction: utils.GetEncodedFilterFunction(filterConfig),
		Region: "usa_east",
		StartRange: startRange,
		EndRange: endRange,
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: map[string]interface{}{
			"url": "webhook", // todo, set up a webhook to receive the data and use HMAC to verify the data
			"compression": "none",
			"max_retry": 5,
			"retry_interval_sec": 1,
			"post_timeout_sec":   10,
		},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("UPDATE_STREAM: error marshaling payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeAPIURL, streamId)

	client, err := getQuickNodeStreamClient(streamId)
	if err != nil {
		return fmt.Errorf("UPDATE_STREAM: error creating client: %w", err)
	}

	// Execute POST request using FastShot
	_, err = client.Build().
		PATCH(url).
		Body().AsJSON(jsonPayload).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("UPDATE_STREAM: error sending request: %w", err)
	}

	return nil
}

func (q *QuickNodeStreamManager) DeleteAddressStream(ctx context.Context, streamId string) error {
	if streamId == "" {
		return fmt.Errorf("DELETE_STREAM: streamId is required")
	}
	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeAPIURL, streamId)
	client, err := getQuickNodeStreamClient(streamId)
	if err != nil {
		return fmt.Errorf("DELETE_STREAM: error creating client: %w", err)
	}

	// Execute DELETE request using FastShot
	_, err = client.Build().
		DELETE(url).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("DELETE_STREAM: error sending request: %w", err)
	}

	return nil
}

func (q *QuickNodeStreamManager) PauseAddressStream(ctx context.Context, streamId string) error {
	if streamId == "" {
		return fmt.Errorf("PAUSE_STREAM: streamId is required")
	}

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeAPIURL, streamId, "pause")
	
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("PAUSE_STREAM: error creating client: %w", err)
	}

	// Execute PATCH request using FastShot
	_, err = client.Build().
		PATCH(url).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("PAUSE_STREAM: error sending request: %w", err)
	}

	return nil
}

func (q *QuickNodeStreamManager) ActivateAddressStream(ctx context.Context, streamId string) error {
	if streamId == "" {
		return fmt.Errorf("ACTIVATE_STREAM: streamId is required")
	}

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeAPIURL, streamId, "activate")
	
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("ACTIVATE_STREAM: error creating client: %w", err)
	}

	// Execute PATCH request using FastShot
	_, err = client.Build().
		PATCH(url).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("ACTIVATE_STREAM: error sending request: %w", err)
	}

	return nil
}

func (q *QuickNodeStreamManager) GetAllStreams(ctx context.Context) ([]*types.StreamReturnPayload, error) {
	url := fmt.Sprintf("%s", streamConf.QuickNodeAPIURL)
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return nil, fmt.Errorf("GET_STREAMS: error creating client: %w", err)
	}

	res, err := client.Build().
		GET(url).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return nil, fmt.Errorf("GET_STREAMS: error sending request: %w", err)
	}
	defer res.RawResponse.Body.Close()

	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
        return nil, fmt.Errorf("GET_STREAMS: unexpected status code: %d", res.StatusCode)
    }

	var response []*types.StreamReturnPayload
	if err := json.NewDecoder(res.RawResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("GET_STREAMS: error decoding response: %w", err)
	}

    return response, nil
}

func (q *QuickNodeStreamManager) createQuickNodeStream(ctx context.Context, payload types.StreamCreationParams, url string) (string, error) {
	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error marshaling payload: %w", err)
	}

	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error creating client: %w", err)
	}

	// Execute POST request using FastShot
	res, err := client.Build().
		POST(url).
		Body().AsJSON(jsonPayload).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return "", fmt.Errorf("CREATE_STREAM: error sending request: %w", err)
	}

	data, err := utils.ParseJSONResponse(res.RawResponse)
	if err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error parsing response: %w", err)
	}

	return data["id"].(string), nil
}

// getConfiguredClient returns a configured FastShot client
func getQuickNodeStreamClient(url string) (*fastshot.ClientBuilder, error) {
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

    nonce := make([]byte, 32)
    _, err := rand.Read(nonce)
    if err != nil {
        return nil, err
    }

    currentTimestamp := time.Now().Unix()

    // create a map to hold the nonce and timestamp
    payloadForHMAC := map[string]interface{}{
        "nonce":     fmt.Sprintf("%x", nonce),
        "timestamp": currentTimestamp,
    }

    client := fastshot.NewClient(url).
        Config().SetTimeout(30 * time.Second).
        Header().Add("accept", "application/json").
        Header().Add("Content-Type", "application/json").
        Header().Add("x-api-key", streamConf.QuickNodeAPIKey).
        Header().Add("X-QN-Nonce", fmt.Sprintf("%x", nonce)).
        Header().Add("X-QN-Timestamp", fmt.Sprintf("%d", currentTimestamp)).
        Header().Add("X-QN-Signature", tokenUtils.GenerateHMACSignature(payloadForHMAC, streamConf.QuickNodePrivateKey))

	return client, nil
}
