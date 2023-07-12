package test

import (
	"context"
	"strings"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	
)

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

// DeployERC20Token deploys an ERC20 token contract using the provided parameters.
// It returns the deployed contract's address, transaction, contract instance, and any error encountered.
// func deployERC20Contract(auth *bind.TransactOpts, client *ethclient.Client, name string, symbol string, decimals uint8, initialSupply *big.Int) (common.Address, *types.Transaction, *contracts.YourERC20Contract, error) {
// 	// address := common.Address{}
// 	// contract := &contracts.YourERC20Contract{}

// 	// // Deploy the ERC20 contract
// 	// address, tx, _, err := contracts.DeployYourERC20Contract(auth, client, name, symbol, decimals, initialSupply)
// 	// if err != nil {
// 	// 	return address, tx, contract, err
// 	// }

// 	// // Retrieve the deployed contract instance
// 	// instance, err := contracts.NewYourERC20Contract(address, client)
// 	// if err != nil {
// 	// 	return address, tx, contract, err
// 	// }

// 	// return address, tx, instance, nil

// }

// deployERC20Contract deploys an ERC20 contract with the provided parameters.
// It returns the address of the deployed contract, the transaction object for the deployment,
// an instance of the deployed contract, and any error that occurred.
func deployTestTokenContract(auth *bind.TransactOpts, client *ethclient.Client, initialSupply *big.Int) (common.Address, error) {
	// Load the compiled ABI of the TestToken contract
	TestTokenABI :=  "{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	testTokenABI, err := abi.JSON(strings.NewReader(TestTokenABI))
	if err != nil {
		return common.Address{}, err
	}

	// Generate the deployment data
	deploymentData, err := testTokenABI.Pack("constructor", initialSupply)
	if err != nil {
		return common.Address{}, err
	}

	// Deploy the TestToken contract
	address, _, _, err := bind.DeployContract(auth, testTokenABI, deploymentData, client)
	if err != nil {
		return common.Address{}, err
	}

	return address, nil
}


// // deployERC20Contract deploys an ERC20 contract with the provided parameters.
// // It returns the address of the deployed contract, the transaction object for the deployment,
// // an instance of the deployed contract, and any error that occurred.
// func deployERC20Contract(auth *bind.TransactOpts, client *ethclient.Client, name string, symbol string, decimals uint8, initialSupply *big.Int) (common.Address, *types.Transaction, *TestToken, error) {
// 	address := common.Address{}
// 	contract := &TestToken{}

// 	// Retrieve the ABI for your ERC20 contract
// 	contractAbi, err := abi.JSON(contract.ABI)
// 	if err != nil {
// 		return address, nil, contract, err
// 	}

// 	// Create the deployment data using the contract ABI and constructor arguments
// 	deploymentData, err := contractAbi.Pack("constructor", name, symbol, decimals, initialSupply)
// 	if err != nil {
// 		return address, nil, contract, err
// 	}

// 	// Deploy the ERC20 contract
// 	address, tx, _, err := bind.DeployContract(auth, contractAbi, deploymentData, client)
// 	if err != nil {
// 		return address, tx, contract, err
// 	}

// 	// Retrieve the deployed contract instance
// 	instance, err := TestToken.NewTest(address, client)
// 	if err != nil {
// 		return address, tx, contract, err
// 	}

// 	return address, tx, instance, nil
// }


