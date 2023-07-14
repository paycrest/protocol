package test

import (
	"context"
	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/paycrest/paycrest-protocol/config"
	TestToken "github.com/paycrest/paycrest-protocol/utils/test/contracts"
)


var serverConf = config.ServerConfig()
// SimulatedBackend extends the backends.SimulatedBackend struct.
type SimulatedBackend struct {
	*backends.SimulatedBackend
}

// SendTransaction sends a transaction to the simulated backend and commits the state.
func (s *SimulatedBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := s.SimulatedBackend.SendTransaction(ctx, tx); err != nil {
		return err
	}
	s.Commit()
	return nil
}

// NewSimulatedBlockchain creates a new instance of SimulatedBackend and returns it.
func NewSimulatedBlockchain() (*SimulatedBackend, error) {
	// Generate a private key for the simulated blockchain
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	// Create a new transactor using the generated private key
	auth := bind.NewKeyedTransactor(privateKey)

	// Set the balance for the transactor's address
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	// Set the genesis account with the transactor's address and balance
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	// Set the block gas limit
	blockGasLimit := uint64(4712388)

	// Create a new simulated backend using the genesis allocation and block gas limit
	client := &SimulatedBackend{backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)}

	return client, nil
}

// deployERC20Contract deploys an ERC20 contract with the provided parameters.
// It returns the address of the deployed contract, the transaction object for the deployment,
// an instance of the deployed contract, and any error that occurred.
func DeployERC20Contract() (bool, error) {
	// Create a new simulated blockchain
	backend, err := NewSimulatedBlockchain()
	if err != nil {
		return false, err
	}


	mnemonic := serverConf.HDWalletMnemonic

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return false, fmt.Errorf("failed to create wallet from mnemonic: %w", err)
	}



	path, err := hdwallet.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return false, fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return false, fmt.Errorf("failed to derive account: %w", err)
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return false, fmt.Errorf("failed to get private key: %w", err)
	}
	publicKey := privateKey.PublicKey


	fromAddress := crypto.PubkeyToAddress(publicKey)
	nonce, err := backend.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return false, err
	}

	gasPrice, err := backend.SuggestGasPrice(context.Background())
	if err != nil {
		return false, err
	}
	// A simulated backend always uses chainID 1337.
	chainID:= backend.Blockchain().Config().ChainID

	//Hard coded method
	//chainID := big.NewInt(1337)
	
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey,chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	initialSupply := new(big.Int)
	initialSupply.SetString("200000000", 10)
	

	address, tx, instance, err := TestToken.DeployTestToken(auth, backend, initialSupply)
	if err != nil {
		return false, err
	}

	backend.Commit()

	fmt.Println(address.Hash())
	fmt.Println(tx.Hash().Hex())

	_ = instance

	return true, nil
}
