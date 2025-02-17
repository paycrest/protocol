package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/aggregator/config"
	"github.com/paycrest/aggregator/ent"
	networkent "github.com/paycrest/aggregator/ent/network"
	"github.com/paycrest/aggregator/services"
	orderService "github.com/paycrest/aggregator/services/order"
	tokenent "github.com/paycrest/aggregator/ent/token"
	"github.com/paycrest/aggregator/storage"
	"github.com/paycrest/aggregator/types"
	"github.com/paycrest/aggregator/utils"
	"github.com/paycrest/aggregator/utils/logger"
)

var serverConf = config.ServerConfig()


// ProviderController is a controller type for provider endpoints
type StreamController struct{}

// NewProviderController creates a new instance of ProviderController with injected services
func NewStreamController() *StreamController {
	return &StreamController{}
}

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

func (ctrl *StreamController) QuicknodeLinkedAddressHook(ctx *gin.Context) {
	// Establish RPC connections
	_, err := setRPCClients(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(string(body))
	var payload map[string]interface{}
    if err := json.Unmarshal(body, &payload); err != nil {
        // http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
	
	transactionHash := payload["TxHash"].(string)
	if transactionHash == "" {
		// lazy return for empty transaction hash
		return
	}

	value := payload["Value"].(string)
	// convert value to big.Int
	bigValue, ok := new(big.Int).SetString(value, 10)
	if !ok {
		// lazy return for invalid value
		return
	}

	identifier := ctx.Request.Header.Get("Stream-Network")
	
	// Create a new TokenTransferEvent
    transferEvent := &types.TokenTransferEvent{
		BlockNumber: uint64(payload["BlockNumber"].(int64)),
        TxHash:      transactionHash,
        From:        payload["from"].(string),
        To:          payload["to"].(string),
		Value:       bigValue,
    }

	// Get the token from the database where the network matches the identifier and token is enabled
	token, err := storage.Client.Token.
		Query().
		Where(
			tokenent.ContractAddressEQ(payload["Token"].(string)),
			tokenent.IsEnabledEQ(true),
		).
		WithNetwork().
		Only(ctx)
	if err != nil {
		logger.Errorf("failed to fetch token: %v", err)
		return
	}

	indexerService := services.NewIndexerService(orderService.NewOrderEVM())
	err = indexerService.StreamERC20Transfer(ctx, rpcClients[identifier], nil, token, transferEvent)
	if err != nil {
		logger.Errorf("failed to stream ERC20 transfer: %v", err)
		return
	}

}
