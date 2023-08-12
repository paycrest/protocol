package database

import (
	"context"
	"flag"
	"fmt"

	"github.com/paycrest/paycrest-protocol/ent/network"
)

func SeedAll() error {
	// Define flags
	seedDB := flag.Bool("seed-db", false, "Seed the database")
	flag.Parse()

	// Run based on flags
	if *seedDB {
		err := SeedDatabase()
		if err != nil {
			return fmt.Errorf("error seeding database: %w", err)
		}
	}

	return nil
}

func SeedDatabase() error {
	client := GetClient()

	// Seed Network
	network, err := client.Network.
		Create().
		SetIdentifier(network.IdentifierPolygonMumbai).
		SetChainID(80001).
		SetRPCEndpoint("https://polygon-mumbai.infura.io/v3/4458cf4d1689497b9a38b1d6bbf05e78").
		SetIsTestnet(true).
		Save(context.Background())
	if err != nil {
		return fmt.Errorf("failed seeding network: %w", err)
	}

	// Seed Token
	_, err = client.Token.
		Create().
		SetSymbol("DERC20").
		SetContractAddress("0xfe4F5145f6e09952a5ba9e956ED0C25e3Fa4c7F1").
		SetDecimals(18).
		SetNetwork(network).
		SetIsEnabled(true).
		Save(context.Background())
	if err != nil {
		return fmt.Errorf("failed seeding token: %w", err)
	}

	return nil
}
