package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/services/contracts"
	"github.com/paycrest/paycrest-protocol/types"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
)

// ReceiveAddressService provides functionality related to managing receive addresses
type ReceiveAddressService struct{}

// NewReceiveAddressService creates a new instance of ReceiveAddressService.
func NewReceiveAddressService() *ReceiveAddressService {
	return &ReceiveAddressService{}
}

// CreateSmartAccount function generates and saves a new EIP-4337 smart contract account address
func (s *ReceiveAddressService) CreateSmartAccount(ctx context.Context, client types.RPCClient, factory *common.Address) (*ent.ReceiveAddress, error) {

	// Connect to RPC endpoint
	var err error
	if client == nil {
		client, err = types.NewEthClient("https://mainnet.infura.io/v3/4818dbcee84d4651a832894818bd4534")
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Initialize contract factory
	if factory == nil {
		// https://github.com/eth-infinitism/account-abstraction/blob/develop/contracts/samples/SimpleAccountFactory.sol
		factoryAddress := common.HexToAddress("0x9406Cc6185a346906296840746125a0E44976454")
		factory = &factoryAddress
	}

	simpleAccountFactory, err := contracts.NewSimpleAccountFactory(*factory, client.(bind.ContractBackend))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize factory contract: %w", err)
	}

	// Get master private key
	ownerAddress, ownerPrivateKeyHex, _ := cryptoUtils.GenerateAccountFromIndex(0)

	// Decode private key
	ownerPrivateKeyBytes, err := hexutil.Decode(ownerPrivateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Use SHA-256 to generate the salt
	hash := sha256.Sum256(ownerPrivateKeyBytes)

	// Create a new big.Int from the hash
	salt := new(big.Int).SetBytes(hash[:])

	// Generate address
	_, err = simpleAccountFactory.GetAddress(nil, common.HexToAddress(ownerAddress), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	// Save address in db
	receiveAddress, err := db.Client.ReceiveAddress.
		Create().
		SetAddress("0xF6F6407410235202CA5Bfa68286a3bBe01F8E5E0").
		SetStatus(receiveaddress.StatusUnused).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save address: %w", err)
	}

	return receiveAddress, nil
}
