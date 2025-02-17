package services

import (
	// "bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/aggregator/config"
	"github.com/paycrest/aggregator/ent"
	networkent "github.com/paycrest/aggregator/ent/network"
	"github.com/paycrest/aggregator/services/contracts"
	"github.com/paycrest/aggregator/storage"
	"github.com/paycrest/aggregator/types"
	"github.com/paycrest/aggregator/utils"
	"github.com/paycrest/aggregator/utils/logger"
	tokenUtils "github.com/paycrest/aggregator/utils/token"
	// tokenUtils "github.com/paycrest/aggregator/utils/token"
	// tokenUtils "github.com/paycrest/aggregator/utils/token"
)

var streamConf = config.StreamConfig()
var rpcClients = map[string]types.RPCClient{}
var linkedAddressStreamName = "paycrest-linked-address-stream" // Format that quicknode uses

// setRPCClients connects to the RPC endpoints of all networks
func setRPCClients(ctx context.Context) ([]*ent.Network, error) {
	isTestnet := false
	if serverConf.Environment != "production" {
		isTestnet = true
	}

	networks, err := storage.Client.Network.
		Query().
		Where(networkent.IsTestnetEQ(isTestnet)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("setRPCClients.fetchNetworks: %w", err)
	}

	// Connect to RPC endpoint
	var client types.RPCClient
	for _, network := range networks {
		if rpcClients[network.Identifier] == nil && !strings.HasPrefix(network.Identifier, "tron") {
			retryErr := utils.Retry(3, 1*time.Second, func() error {
				client, err = types.NewEthClient(network.RPCEndpoint)
				return err
			})
			if retryErr != nil {
				logger.Errorf("failed to connect to %s RPC %v", network.Identifier, retryErr)
				continue
			}

			rpcClients[network.Identifier] = client
		}
	}

	return networks, nil
}

// QuickNodeStreamManager is a service for managing streams using QuickNode
type QuickNodeStreamManager struct{}

func NewQuickNodeStreamManager() *QuickNodeStreamManager{
	return &QuickNodeStreamManager{}
}

// This function will need to be called everytime a new receive address is generated
func (q *QuickNodeStreamManager) CreateAddressStream(ctx context.Context, order *ent.PaymentOrder, token *ent.Token, identifier string, startRange int) (string, error) {
	_, err := setRPCClients(ctx)
	if err != nil {
        logger.Errorf("CreateAddressStream: error setting RPC clients: %v", err)
        return "", fmt.Errorf("CreateAddressStream: error setting RPC clients: %w", err)
    }

	client := rpcClients[identifier]
	// Connect to RPC endpoint
	if client == nil {
		retryErr := utils.Retry(3, 1*time.Second, func() error {
			client, err = types.NewEthClient(token.Edges.Network.RPCEndpoint)
			return err
		})
		if retryErr != nil {
			logger.Errorf("CreateAddressStream: error connecting to RPC endpoint: %v", retryErr)
			return "", retryErr
		}
	}
	var addressToWatch string

	if order != nil {
		token = order.Edges.Token
		addressToWatch = order.Edges.ReceiveAddress.Address
	}

	// Fetch current block header
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
        logger.Errorf("CreateAddressStream: error fetching block header: %v", err)
        return "", err
    }
	endRange := header.Number.Int64()

	addresses := []common.Address{}
	if addressToWatch != "" {
		fromRange := int(5000)
		if token.Edges.Network.Identifier == "bnb-smart-chain" {
			startRange = 1000
		}
		addresses = []common.Address{common.HexToAddress(addressToWatch)}
		startRange = int(endRange) - fromRange
	} else {
		// addresses = []common.Address{common.HexToAddress("0x0f13743a0C27D898EC62F3444Ed6a99E009a8912")} // serve for test purpose
		startRange = int(endRange) - 100
	}

	// filterConfig := types.FilterConfig{
	// 	Addresses: addresses,
	// 	ERC20Tokens: []common.Address{common.HexToAddress("0xaf88d065e77c8cC2239327C5EDb3A432268e5831")},
	// 	Abi: contracts.TransferEventABI,
	// }

	// length of addresses should be at least 1
	if len(addresses) == 0 {
		return "", nil
	}

	filterConfig := types.FilterConfig{
		Addresses: addresses,
		ERC20Tokens: []common.Address{common.HexToAddress(token.ContractAddress)},
		Abi: contracts.ERC20TokenABI,
	}

	allFunctions, err := q.GetAllFunctions(ctx)
	if err != nil {
		logger.Errorf("CreateAddressStream: error fetching all functions: %v", err)
		return "", err
	}

	var filterFunctionExists bool
	var code string

	// Check if the function already exists
	for _, function := range allFunctions {
		if function["name"].(string) == linkedAddressStreamName {
			filterFunctionExists = true
			code = function["code"].(string)
			break
		}
	}

	var quickNodeFilterFunctionCodePayload = types.QuickNodeFunctionPayload{
		Name: linkedAddressStreamName,
		Description: "Paycrest Filter function for linked addresses",
		Kind: "nodejs-qn:20",
		Code: utils.GetEncodedFilterFunction(filterConfig),
		Binary: false,
        Limits: map[string]int{
            "timeout": 5000,
        },
	}

	if !filterFunctionExists {
		code, err = CreateQuickNodeFunction(quickNodeFilterFunctionCodePayload)
		if err != nil {
			logger.Errorf("CreateAddressStream: error creating new function: %v", err)
			return "", err
		}
	}

	getAllStreams, err := q.GetAllStreams(ctx)
	if err != nil {
        logger.Errorf("CreateAddressStream: error fetching all streams: %v", err)
        return "", err
    }

	for _, stream := range getAllStreams {
		if stream.Name == linkedAddressStreamName && stream.Network == identifier {
			// Stream already exists
			// update the stream with the new range
			err := q.UpdateAddressStream(ctx, stream.ID, startRange, int(endRange), code)
			if err != nil {
                logger.Errorf("CreateAddressStream: error updating existing stream: %v", err)
                return "", err
            }
			return stream.ID, nil
		}
	}

	destinationAttributes, err := generateHeaderForStream()
	if err != nil {
		logger.Errorf("CreateAddressStream: error generating header for stream: %v", err)
		return "", err
	}

	params := types.StreamCreationParams{
		Name: linkedAddressStreamName,
		Network: "arbitrum-mainnet", // identifier,
		Dataset: "block_with_receipts",
		FilterFunction: code,
		Region: "usa_east",
		// StartRange:            305080560, // 305064772
		// EndRange:              305080570,
		StartRange: startRange,
		EndRange: int(endRange),
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		IncludeStreamMetadata: "body",
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: destinationAttributes,
	}

	streamId, err := q.createQuickNodeStream(ctx, params, streamConf.QuickNodeStreamAPIURL)
	if err != nil {
        logger.Errorf("CreateAddressStream: error creating new stream: %v", err)
        return "", err
    }

    return streamId, nil
}

func (q *QuickNodeStreamManager) UpdateAddressStream(ctx context.Context, streamId string, startRange int, endRange int, code string) error {
	if endRange <= startRange {
		return fmt.Errorf("UPDATE_STREAM: endRange must be greater than startRange")
	}
	if streamId == "" {
		return fmt.Errorf("UPDATE_STREAM: streamId is required")
	}

	if code == "" {
		return fmt.Errorf("UPDATE_STREAM: code is required")
	}

	destinationAttributes, err := generateHeaderForStream()
	if err != nil {
		logger.Errorf("CreateAddressStream: error generating header for stream: %v", err)
		return err
	}

	params := types.StreamCreationParams{
		Name: "PaycrestAddressStream",
		FilterFunction: code,
		Region: "usa_east",
		StartRange: startRange,
		EndRange: endRange,
		// StartRange:            305080560, // 305064772
		// EndRange:              305080570,
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: destinationAttributes,
	}

	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeStreamAPIURL, streamId)

	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("UPDATE_STREAM: error creating client: %w", err)
	}

	// Execute POST request using FastShot
	_, err = client.Build().
		PATCH("").
		Body().AsJSON(&params).
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
	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeStreamAPIURL, streamId)
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("DELETE_STREAM: error creating client: %w", err)
	}

	// Execute DELETE request using FastShot
	_, err = client.Build().
		DELETE("").
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

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeStreamAPIURL, streamId, "pause")
	
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("PAUSE_STREAM: error creating client: %w", err)
	}

	// Execute PATCH request using FastShot
	_, err = client.Build().
		PATCH("").
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

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeStreamAPIURL, streamId, "activate")
	
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("ACTIVATE_STREAM: error creating client: %w", err)
	}

	// Execute PATCH request using FastShot
	_, err = client.Build().
		PATCH("").
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("ACTIVATE_STREAM: error sending request: %w", err)
	}

	return nil
}

func (q *QuickNodeStreamManager) GetAllStreams(ctx context.Context) ([]*types.StreamReturnPayload, error) {
	url := streamConf.QuickNodeStreamAPIURL
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return nil, fmt.Errorf("GET_STREAMS: error creating client: %w", err)
	}

	res, err := client.Build().
		GET("").
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return nil, fmt.Errorf("GET_STREAMS: error sending request: %w", err)
	}

	defer res.RawResponse.Body.Close()

	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return nil, fmt.Errorf("GET_STREAMS: unexpected status code: %d", res.StatusCode())
    }

	var response []*types.StreamReturnPayload
	body, _ := io.ReadAll(res.RawResponse.Body)
	json.Unmarshal(body, &response)
    return response, nil
}

func (q *QuickNodeStreamManager) createQuickNodeStream(_ context.Context, payload types.StreamCreationParams, url string) (string, error) {
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error creating client: %w", err)
	}

	res, err := client.Build().
		POST("").
		Body().AsJSON(&payload).
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
		return nil, fmt.Errorf("QuickNode: URL is required")
	}

	headers := map[string]string{
		"accept": "application/json",
		"Content-Type": "application/json",
		"x-api-key": streamConf.QuickNodeAPIKey,
	}

    client := fastshot.NewClient(url).
        Config().SetTimeout(30 * time.Second).
		Header().AddAll(headers)
	
	return client, nil
}

func generateHeaderForStream() (types.DestinationAttributes, error) {
	nonce := make([]byte, 32)
    _, err := rand.Read(nonce)
    if err != nil {
        return types.DestinationAttributes{}, err
    }

    currentTimestamp := time.Now().Unix()

    // // create a map to hold the nonce and timestamp
    payloadForHMAC := map[string]interface{}{
        "nonce":     fmt.Sprintf("%x", nonce),
        "timestamp": fmt.Sprintf("%d", currentTimestamp),
    }

	signature := tokenUtils.GenerateHMACSignature(payloadForHMAC, streamConf.QuickNodePrivateKey)

	destinationHeader := types.DestinationAttributes{
			URL: `https://88c5-86-13-108-4.ngrok-free.app/v1/stream/quicknode-linked-addresses-hook`, // Set the correct URL for the webhook
			Compression: "none",
			MaxRetry: 5,
			RetryIntervalSec: 1,
			PostTimeoutSec:   10,
			Headers: types.DestinationAttributesHeaders{
				ClientType: "quicknode",
				Authorization: signature,
				Nonce: fmt.Sprintf("%x", nonce),
				Timestamp: fmt.Sprintf("%d", currentTimestamp),
			},
		}
	return destinationHeader, nil
}

// This function is expected to replace {IndexBlockchainEvents()}
func (q *QuickNodeStreamManager) GetFunctionAddressStreamData(ctx context.Context, streamId string, identifier string, blockNumber int, key string) error {
	if streamId == "" {
		return fmt.Errorf("ACTIVATE_STREAM: streamId is required")
	}
	// /call/result_only=true

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeFunctionAPIURL, streamId, "call")
	
	client, err := getQuickNodeStreamClient(url)
	if err != nil {
		return fmt.Errorf("ACTIVATE_STREAM: error creating client: %w", err)
	}

	payload := map[string]interface{}{
		"network": identifier,
		"dataset": "block_with_receipts",
		"block_number": blockNumber,
		"result_only": true,
		"user_data": map[string]interface{}{
			"key": key,
		},
	}

	// Execute PATCH request using FastShot
	res, err := client.Build().
		POST("").
		Body().AsJSON(&payload).
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return fmt.Errorf("ACTIVATE_STREAM: error sending request: %w", err)
	}

	defer res.RawResponse.Body.Close()

	// @todo the response doesn allow to pull data in the format we want, something I need to work on

	return nil
}

// @dev: This function is used to create a new function on QuickNode, so always check if the function already exists before creating a new one
func CreateQuickNodeFunction(functionData types.QuickNodeFunctionPayload) (string, error) {
	client, err := getQuickNodeStreamClient(streamConf.QuickNodeFunctionAPIURL)
	if err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error creating client: %w", err)
	}

	res, err := client.Build().
		POST("").
		Body().AsJSON(&functionData).
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

	var execPayload map[string]interface{}
	if err := json.Unmarshal([]byte(data["exec"].(string)), &execPayload); err != nil {
		return "", fmt.Errorf("CREATE_STREAM: error parsing response payload: %w", err)
	}

	return execPayload["code"].(string), nil
}

func (q *QuickNodeStreamManager) GetAllFunctions(ctx context.Context) ([]map[string]interface{}, error) {
	client, err := getQuickNodeStreamClient(streamConf.QuickNodeFunctionAPIURL)
	if err != nil {
		return nil, fmt.Errorf("GET_FILTERFUNCTION: error creating client: %w", err)
	}

	res, err := client.Build().
		GET("").
		Retry().Set(3, 1*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error sending request: %v", err)
		return nil, fmt.Errorf("GET_FILTERFUNCTION: error sending request: %w", err)
	}

	defer res.RawResponse.Body.Close()

	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return nil, fmt.Errorf("GET_FILTERFUNCTION: unexpected status code: %d", res.StatusCode())
    }

	var response map[string]interface{}
	body, _ := io.ReadAll(res.RawResponse.Body)
	json.Unmarshal(body, &response)

	datas, ok := response["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("GET_FILTERFUNCTION: error parsing response")
	}

	result := make([]map[string]interface{}, len(datas))
	for i, data := range datas {
		result[i] = data.(map[string]interface{})
	}
	return result, nil
}