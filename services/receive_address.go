package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/paycrest/paycrest-protocol/config"
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

	// Get master account
	ownerAddress, _, _ := cryptoUtils.GenerateAccountFromIndex(0)

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	// Create a new big.Int from the hash
	salt := new(big.Int).SetBytes(nonce)

	// Generate address
	address, err := simpleAccountFactory.GetAddress(nil, common.HexToAddress(ownerAddress), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	saltEncrypted, err := cryptoUtils.EncryptPlain([]byte(salt.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt salt: %w", err)
	}

	// Save address in db
	conf := config.OrderConfig()
	receiveAddress, err := db.Client.ReceiveAddress.
		Create().
		SetAddress(address.Hex()).
		SetSalt(saltEncrypted).
		SetStatus(receiveaddress.StatusUnused).
		SetValidUntil(time.Now().Add(conf.ReceiveAddressValidity)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save address: %w", err)
	}

	return receiveAddress, nil
}
