// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TestTokenMetaData contains all meta data concerning the TestToken contract.
var TestTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b50604051620018dd380380620018dd83398181016040528101906200003691906200028c565b6040518060400160405280600981526020017f54657374546f6b656e00000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f544f5400000000000000000000000000000000000000000000000000000000008152508160039081620000b3919062000517565b508060049081620000c5919062000517565b505050620000da3382620000e160201b60201c565b506200070c565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000152576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620001499062000659565b60405180910390fd5b620001655f83836200024660201b60201c565b8060025f828254620001789190620006a6565b92505081905550805f808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508173ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620002279190620006f1565b60405180910390a3620002425f83836200024b60201b60201c565b5050565b505050565b505050565b5f80fd5b5f819050919050565b620002688162000254565b811462000273575f80fd5b50565b5f8151905062000286816200025d565b92915050565b5f60208284031215620002a457620002a362000250565b5b5f620002b38482850162000276565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200033857607f821691505b6020821081036200034e576200034d620002f3565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620003b27fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000375565b620003be868362000375565b95508019841693508086168417925050509392505050565b5f819050919050565b5f620003ff620003f9620003f38462000254565b620003d6565b62000254565b9050919050565b5f819050919050565b6200041a83620003df565b62000432620004298262000406565b84845462000381565b825550505050565b5f90565b620004486200043a565b620004558184846200040f565b505050565b5b818110156200047c57620004705f826200043e565b6001810190506200045b565b5050565b601f821115620004cb57620004958162000354565b620004a08462000366565b81016020851015620004b0578190505b620004c8620004bf8562000366565b8301826200045a565b50505b505050565b5f82821c905092915050565b5f620004ed5f1984600802620004d0565b1980831691505092915050565b5f620005078383620004dc565b9150826002028217905092915050565b6200052282620002bc565b67ffffffffffffffff8111156200053e576200053d620002c6565b5b6200054a825462000320565b6200055782828562000480565b5f60209050601f8311600181146200058d575f841562000578578287015190505b620005848582620004fa565b865550620005f3565b601f1984166200059d8662000354565b5f5b82811015620005c6578489015182556001820191506020850194506020810190506200059f565b86831015620005e65784890151620005e2601f891682620004dc565b8355505b6001600288020188555050505b505050505050565b5f82825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f2061646472657373005f82015250565b5f62000641601f83620005fb565b91506200064e826200060b565b602082019050919050565b5f6020820190508181035f830152620006728162000633565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f620006b28262000254565b9150620006bf8362000254565b9250828201905080821115620006da57620006d962000679565b5b92915050565b620006eb8162000254565b82525050565b5f602082019050620007065f830184620006e0565b92915050565b6111c3806200071a5f395ff3fe608060405234801561000f575f80fd5b50600436106100a7575f3560e01c8063395093511161006f578063395093511461016557806370a082311461019557806395d89b41146101c5578063a457c2d7146101e3578063a9059cbb14610213578063dd62ed3e14610243576100a7565b806306fdde03146100ab578063095ea7b3146100c957806318160ddd146100f957806323b872dd14610117578063313ce56714610147575b5f80fd5b6100b3610273565b6040516100c09190610add565b60405180910390f35b6100e360048036038101906100de9190610b8e565b610303565b6040516100f09190610be6565b60405180910390f35b610101610325565b60405161010e9190610c0e565b60405180910390f35b610131600480360381019061012c9190610c27565b61032e565b60405161013e9190610be6565b60405180910390f35b61014f61035c565b60405161015c9190610c92565b60405180910390f35b61017f600480360381019061017a9190610b8e565b610364565b60405161018c9190610be6565b60405180910390f35b6101af60048036038101906101aa9190610cab565b61039a565b6040516101bc9190610c0e565b60405180910390f35b6101cd6103df565b6040516101da9190610add565b60405180910390f35b6101fd60048036038101906101f89190610b8e565b61046f565b60405161020a9190610be6565b60405180910390f35b61022d60048036038101906102289190610b8e565b6104e4565b60405161023a9190610be6565b60405180910390f35b61025d60048036038101906102589190610cd6565b610506565b60405161026a9190610c0e565b60405180910390f35b60606003805461028290610d41565b80601f01602080910402602001604051908101604052809291908181526020018280546102ae90610d41565b80156102f95780601f106102d0576101008083540402835291602001916102f9565b820191905f5260205f20905b8154815290600101906020018083116102dc57829003601f168201915b5050505050905090565b5f8061030d610588565b905061031a81858561058f565b600191505092915050565b5f600254905090565b5f80610338610588565b9050610345858285610752565b6103508585856107dd565b60019150509392505050565b5f6012905090565b5f8061036e610588565b905061038f8185856103808589610506565b61038a9190610d9e565b61058f565b600191505092915050565b5f805f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6060600480546103ee90610d41565b80601f016020809104026020016040519081016040528092919081815260200182805461041a90610d41565b80156104655780601f1061043c57610100808354040283529160200191610465565b820191905f5260205f20905b81548152906001019060200180831161044857829003601f168201915b5050505050905090565b5f80610479610588565b90505f6104868286610506565b9050838110156104cb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104c290610e41565b60405180910390fd5b6104d8828686840361058f565b60019250505092915050565b5f806104ee610588565b90506104fb8185856107dd565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f33905090565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036105fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105f490610ecf565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361066b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161066290610f5d565b60405180910390fd5b8060015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107459190610c0e565b60405180910390a3505050565b5f61075d8484610506565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107d757818110156107c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107c090610fc5565b60405180910390fd5b6107d6848484840361058f565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361084b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161084290611053565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108b0906110e1565b60405180910390fd5b6108c4838383610a49565b5f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610947576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161093e9061116f565b60405180910390fd5b8181035f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550815f808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a309190610c0e565b60405180910390a3610a43848484610a4e565b50505050565b505050565b505050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015610a8a578082015181840152602081019050610a6f565b5f8484015250505050565b5f601f19601f8301169050919050565b5f610aaf82610a53565b610ab98185610a5d565b9350610ac9818560208601610a6d565b610ad281610a95565b840191505092915050565b5f6020820190508181035f830152610af58184610aa5565b905092915050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610b2a82610b01565b9050919050565b610b3a81610b20565b8114610b44575f80fd5b50565b5f81359050610b5581610b31565b92915050565b5f819050919050565b610b6d81610b5b565b8114610b77575f80fd5b50565b5f81359050610b8881610b64565b92915050565b5f8060408385031215610ba457610ba3610afd565b5b5f610bb185828601610b47565b9250506020610bc285828601610b7a565b9150509250929050565b5f8115159050919050565b610be081610bcc565b82525050565b5f602082019050610bf95f830184610bd7565b92915050565b610c0881610b5b565b82525050565b5f602082019050610c215f830184610bff565b92915050565b5f805f60608486031215610c3e57610c3d610afd565b5b5f610c4b86828701610b47565b9350506020610c5c86828701610b47565b9250506040610c6d86828701610b7a565b9150509250925092565b5f60ff82169050919050565b610c8c81610c77565b82525050565b5f602082019050610ca55f830184610c83565b92915050565b5f60208284031215610cc057610cbf610afd565b5b5f610ccd84828501610b47565b91505092915050565b5f8060408385031215610cec57610ceb610afd565b5b5f610cf985828601610b47565b9250506020610d0a85828601610b47565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610d5857607f821691505b602082108103610d6b57610d6a610d14565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610da882610b5b565b9150610db383610b5b565b9250828201905080821115610dcb57610dca610d71565b5b92915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f775f8201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b5f610e2b602583610a5d565b9150610e3682610dd1565b604082019050919050565b5f6020820190508181035f830152610e5881610e1f565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f206164645f8201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b5f610eb9602483610a5d565b9150610ec482610e5f565b604082019050919050565b5f6020820190508181035f830152610ee681610ead565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f2061646472655f8201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b5f610f47602283610a5d565b9150610f5282610eed565b604082019050919050565b5f6020820190508181035f830152610f7481610f3b565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e63650000005f82015250565b5f610faf601d83610a5d565b9150610fba82610f7b565b602082019050919050565b5f6020820190508181035f830152610fdc81610fa3565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f2061645f8201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b5f61103d602583610a5d565b915061104882610fe3565b604082019050919050565b5f6020820190508181035f83015261106a81611031565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f20616464725f8201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b5f6110cb602383610a5d565b91506110d682611071565b604082019050919050565b5f6020820190508181035f8301526110f8816110bf565b9050919050565b7f45524332303a207472616e7366657220616d6f756e74206578636565647320625f8201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b5f611159602683610a5d565b9150611164826110ff565b604082019050919050565b5f6020820190508181035f8301526111868161114d565b905091905056fea2646970667358221220b9aacc299b85dbeb860cc735af32e70cec10e44792be4124155e61e6d346ce4764736f6c63430008140033",
}

// TestTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use TestTokenMetaData.ABI instead.
var TestTokenABI = TestTokenMetaData.ABI

// TestTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestTokenMetaData.Bin instead.
var TestTokenBin = TestTokenMetaData.Bin

// DeployTestToken deploys a new Ethereum contract, binding an instance of TestToken to it.
func DeployTestToken(auth *bind.TransactOpts, backend bind.ContractBackend, initialSupply *big.Int) (common.Address, *types.Transaction, *TestToken, error) {
	parsed, err := TestTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestTokenBin), backend, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestToken{TestTokenCaller: TestTokenCaller{contract: contract}, TestTokenTransactor: TestTokenTransactor{contract: contract}, TestTokenFilterer: TestTokenFilterer{contract: contract}}, nil
}

// TestToken is an auto generated Go binding around an Ethereum contract.
type TestToken struct {
	TestTokenCaller     // Read-only binding to the contract
	TestTokenTransactor // Write-only binding to the contract
	TestTokenFilterer   // Log filterer for contract events
}

// TestTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestTokenSession struct {
	Contract     *TestToken        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestTokenCallerSession struct {
	Contract *TestTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestTokenTransactorSession struct {
	Contract     *TestTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestTokenRaw struct {
	Contract *TestToken // Generic contract binding to access the raw methods on
}

// TestTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestTokenCallerRaw struct {
	Contract *TestTokenCaller // Generic read-only contract binding to access the raw methods on
}

// TestTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestTokenTransactorRaw struct {
	Contract *TestTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestToken creates a new instance of TestToken, bound to a specific deployed contract.
func NewTestToken(address common.Address, backend bind.ContractBackend) (*TestToken, error) {
	contract, err := bindTestToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestToken{TestTokenCaller: TestTokenCaller{contract: contract}, TestTokenTransactor: TestTokenTransactor{contract: contract}, TestTokenFilterer: TestTokenFilterer{contract: contract}}, nil
}

// NewTestTokenCaller creates a new read-only instance of TestToken, bound to a specific deployed contract.
func NewTestTokenCaller(address common.Address, caller bind.ContractCaller) (*TestTokenCaller, error) {
	contract, err := bindTestToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestTokenCaller{contract: contract}, nil
}

// NewTestTokenTransactor creates a new write-only instance of TestToken, bound to a specific deployed contract.
func NewTestTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TestTokenTransactor, error) {
	contract, err := bindTestToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestTokenTransactor{contract: contract}, nil
}

// NewTestTokenFilterer creates a new log filterer instance of TestToken, bound to a specific deployed contract.
func NewTestTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TestTokenFilterer, error) {
	contract, err := bindTestToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestTokenFilterer{contract: contract}, nil
}

// bindTestToken binds a generic wrapper to an already deployed contract.
func bindTestToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken *TestTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken.Contract.TestTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken *TestTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken.Contract.TestTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken *TestTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken.Contract.TestTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken *TestTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken *TestTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken *TestTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestToken *TestTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestToken *TestTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestToken.Contract.Allowance(&_TestToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestToken *TestTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestToken.Contract.Allowance(&_TestToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestToken *TestTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestToken *TestTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestToken.Contract.BalanceOf(&_TestToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestToken *TestTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestToken.Contract.BalanceOf(&_TestToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestToken *TestTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestToken *TestTokenSession) Decimals() (uint8, error) {
	return _TestToken.Contract.Decimals(&_TestToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestToken *TestTokenCallerSession) Decimals() (uint8, error) {
	return _TestToken.Contract.Decimals(&_TestToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken *TestTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken *TestTokenSession) Name() (string, error) {
	return _TestToken.Contract.Name(&_TestToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken *TestTokenCallerSession) Name() (string, error) {
	return _TestToken.Contract.Name(&_TestToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken *TestTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken *TestTokenSession) Symbol() (string, error) {
	return _TestToken.Contract.Symbol(&_TestToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken *TestTokenCallerSession) Symbol() (string, error) {
	return _TestToken.Contract.Symbol(&_TestToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestToken *TestTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestToken *TestTokenSession) TotalSupply() (*big.Int, error) {
	return _TestToken.Contract.TotalSupply(&_TestToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestToken *TestTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _TestToken.Contract.TotalSupply(&_TestToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestToken *TestTokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.Approve(&_TestToken.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.Approve(&_TestToken.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestToken *TestTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestToken *TestTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.DecreaseAllowance(&_TestToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestToken *TestTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.DecreaseAllowance(&_TestToken.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestToken *TestTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestToken *TestTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.IncreaseAllowance(&_TestToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestToken *TestTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.IncreaseAllowance(&_TestToken.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.Transfer(&_TestToken.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.Transfer(&_TestToken.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.TransferFrom(&_TestToken.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TestToken *TestTokenTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken.Contract.TransferFrom(&_TestToken.TransactOpts, from, to, amount)
}

// TestTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestToken contract.
type TestTokenApprovalIterator struct {
	Event *TestTokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestTokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestTokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestTokenApproval represents a Approval event raised by the TestToken contract.
type TestTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestToken *TestTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TestTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TestTokenApprovalIterator{contract: _TestToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestToken *TestTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestTokenApproval)
				if err := _TestToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestToken *TestTokenFilterer) ParseApproval(log types.Log) (*TestTokenApproval, error) {
	event := new(TestTokenApproval)
	if err := _TestToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestToken contract.
type TestTokenTransferIterator struct {
	Event *TestTokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestTokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestTokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestTokenTransfer represents a Transfer event raised by the TestToken contract.
type TestTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestToken *TestTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestTokenTransferIterator{contract: _TestToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestToken *TestTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestTokenTransfer)
				if err := _TestToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestToken *TestTokenFilterer) ParseTransfer(log types.Log) (*TestTokenTransfer, error) {
	event := new(TestTokenTransfer)
	if err := _TestToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
