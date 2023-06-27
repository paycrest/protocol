package services

import (
	"context"
	"fmt"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
)

type ReceiveAddressService struct {
	db *ent.Client
}

func NewReceiveAddressService(db *ent.Client) *ReceiveAddressService {
	return &ReceiveAddressService{
		db: db,
	}
}

func (s *ReceiveAddressService) GenerateAndSaveAddress(ctx context.Context) (string, error) {
	// TODO: query db to know number of receive addresses
	// accountIndex should be num_of_address + 1
	accountIndex := 0

	address, _, err := crypto.GenerateReceiveAddress(accountIndex)

	// TODO: Save the address and privateKey to the database using ent
	// Save the address and privateKey to the database using ent
	// statuses: unused, partial, used, expired
	_, err = s.db.ReceiveAddress.
		Create().
		SetAddress(address).
		SetAccountIndex(accountIndex).
		SetStatus("active").
		Save(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to save address: %w", err)
	}

	return address, nil
}
