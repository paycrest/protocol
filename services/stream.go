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
		startRange = int(endRange) - 100
	}

	filterConfig := types.FilterConfig{
		Addresses: addresses,
		ERC20Tokens: []common.Address{common.HexToAddress(token.ContractAddress)},
		Abi: contracts.ERC20TokenABI,
		ListName: "PaycrestLinkedAddresses",
	}

	getAllStreams, err := q.GetAllStreams(ctx)
	if err != nil {
        logger.Errorf("CreateAddressStream: error fetching all streams: %v", err)
        return "", err
    }

	for _, stream := range getAllStreams {
		if stream.Name == "PaycrestLinkedAddressStream" && stream.Network == identifier {
			// Stream already exists
			// update the stream with the new range
			err := q.UpdateAddressStream(ctx, stream.ID, startRange, int(endRange), filterConfig)
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

	fmt.Println("..........y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.y.yy..destinationAttributes", destinationAttributes)

	params := types.StreamCreationParams{
		Name: "PaycrestLinkedAddressStream",
		Network: "arbitrum-mainnet",
		Dataset: "block_with_receipts",
		// FilterFunction: utils.GetEncodedFilterFunction(filterConfig),
		Region: "usa_east",
		StartRange: startRange,
		EndRange: int(endRange),
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		IncludeStreamMetadata: "body",
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: destinationAttributes,
	}

	streamId, err := q.createQuickNodeStream(ctx, params, streamConf.QuickNodeAPIURL)
	if err != nil {
        logger.Errorf("CreateAddressStream: error creating new stream: %v", err)
        return "", err
    }

	fmt.Println(".0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0..0.0.streamId", streamId)
    return streamId, nil
}

func (q *QuickNodeStreamManager) UpdateAddressStream(ctx context.Context, streamId string, startRange int, endRange int, filterConfig types.FilterConfig) error {
	if endRange <= startRange {
		return fmt.Errorf("UPDATE_STREAM: endRange must be greater than startRange")
	}
	if streamId == "" {
		return fmt.Errorf("UPDATE_STREAM: streamId is required")
	}

	destinationAttributes, err := generateHeaderForStream()
	if err != nil {
		logger.Errorf("CreateAddressStream: error generating header for stream: %v", err)
		return err
	}

	params := types.StreamCreationParams{
		Name: "PaycrestAddressStream",
		// FilterFunction: utils.GetEncodedFilterFunction(filterConfig),
		Region: "usa_east",
		StartRange: startRange,
		EndRange: endRange,
		DatasetBatchSize: 2, // need to be verify to know how many transactions to fetch at a time @chibie
		Status: "active",
		Destination: "webhook",
		DestinationAttributes: destinationAttributes,
	}

	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeAPIURL, streamId)

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
	url := fmt.Sprintf("%s/%s", streamConf.QuickNodeAPIURL, streamId)
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

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeAPIURL, streamId, "pause")
	
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

	url := fmt.Sprintf("%s/%s/%s", streamConf.QuickNodeAPIURL, streamId, "activate")
	
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
	url := streamConf.QuickNodeAPIURL
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

	fmt.Println("res.040404040404040404004040..4.4.4.4.4.4.4.4.4.4.4.4.4.4.4RawResponse", res.RawResponse)

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

	fmt.Println("..........url is............................................")

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
			URL: `/v1/stream/quicknode-linked-addresses-hook`, // Set the correct URL for the webhook
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
