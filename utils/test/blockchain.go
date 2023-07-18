package test

import (
	"context"
	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

    "golang.org/x/crypto/sha3"
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

// func FundAddressWithTestToken(address, amount) {
	
// }


func FundAddressWithTestToken(address common.Address, amount *big.Int) error {
	// Create a new simulated blockchain
	backend, err := NewSimulatedBlockchain()
	if err != nil {
		return err
	}

	// Fetch the mnemonic for the HD wallet from server configuration
	mnemonic := serverConf.HDWalletMnemonic

	// Create a new wallet using the mnemonic
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to create wallet from mnemonic: %w", err)
	}

	// Define the derivation path for the Ethereum account
	path, err := hdwallet.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return fmt.Errorf("failed to parse derivation path: %w", err)
	}

	// Derive an Ethereum account from the wallet using the derivation path
	account, err := wallet.Derive(path, false)
	if err != nil {
		return fmt.Errorf("failed to derive account: %w", err)
	}

	// Get the private key associated with the derived account
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	// Derive the public key from the private key
	publicKey := privateKey.PublicKey

	// Get the sender address from the public key
	fromAddress := crypto.PubkeyToAddress(publicKey)

	// Fetch the next available nonce for the sender address
	nonce, err := backend.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to fetch nonce: %w", err)
	}

	value := big.NewInt(0) // Transaction value in wei (0 eth)

	// Suggest gas price for the transaction
	gasPrice, err := backend.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %w", err)
	}

	// Hardcoded recipient and ERC20 token contract addresses
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// Generate the method ID for the ERC20 `transfer` function
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// Print the generated method ID (4 bytes)
	fmt.Println(hexutil.Encode(methodID))

	// Pad the recipient address to 32 bytes
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))

	// Pad the token amount to 32 bytes
	amountPadded := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, amountPadded...)

	// Estimate the gas required for the transaction
	gasLimit, err := backend.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Create a new transaction with the specified details
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// Set the chain ID to 1337 for EIP-155
	chainID := big.NewInt(1337)

	// Sign the transaction using the chain ID as the EIP-155 chain ID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the signed transaction to the simulated backend
	err = backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// Print the hash of the sent transaction
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	return nil
}


//return private and public key