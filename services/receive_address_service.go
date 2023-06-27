package services

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	db "github.com/paycrest/paycrest-protocol/database"
)

type ReceiveAddressService struct {
	wallet       *hdwallet.Wallet
	initialIndex int
}

func NewReceiveAddressService(wallet *hdwallet.Wallet, initialIndex int) *ReceiveAddressService {
	return &ReceiveAddressService{
		wallet:       wallet,
		initialIndex: initialIndex,
	}
}

func (s *ReceiveAddressService) GenerateAndSaveAddress() (string, string, error) {
	path, err := hdwallet.ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", s.initialIndex))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := s.wallet.Derive(path, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive account: %w", err)
	}

	privateKey, err := s.wallet.PrivateKey(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to get private key: %w", err)
	}

	privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))
	address := account.Address.Hex()

	// TODO: Save the address and privateKey to the database using ent
		// Save the address and privateKey to the database using ent
		receivedAddress, err := s.db.ReceiveAddress.
		Create().
		SetAddress(address).
		SetPrivateKey(privateKeyHex).
		SetStatus("active").
		Save(s.ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to save address: %w", err)
	}

	return address, privateKeyHex, nil
}
