package storage

import (
	"context"
	"flag"
	"fmt"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/fiatcurrency"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/shopspring/decimal"
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
	ctx := context.Background()

	// Seed Network
	network, err := client.Network.
		Create().
		SetIdentifier("polygon-mumbai").
		SetChainID(80001).
		SetRPCEndpoint("wss://polygon-mumbai.infura.io/ws/v3/4458cf4d1689497b9a38b1d6bbf05e78").
		SetIsTestnet(true).
		Save(ctx)
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
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed seeding token: %w", err)
	}

	// Seed Provision Buckets
	for _, currency := range []string{"NGN", "KES"} {
		sampleBuckets := make([]*ent.ProvisionBucketCreate, 0, 6)

		// Add sample provision buckets.
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(20000001.00)).
			SetMaxAmount(decimal.NewFromFloat(100000000.00)).
			SetCurrency(currency),
		)
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(5000001.00)).
			SetMaxAmount(decimal.NewFromFloat(20000000.00)).
			SetCurrency(currency),
		)
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(2000001.00)).
			SetMaxAmount(decimal.NewFromFloat(5000000.00)).
			SetCurrency(currency),
		)
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(500001.00)).
			SetMaxAmount(decimal.NewFromFloat(2000000.00)).
			SetCurrency(currency),
		)
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(50001.00)).
			SetMaxAmount(decimal.NewFromFloat(500000.00)).
			SetCurrency(currency),
		)
		sampleBuckets = append(sampleBuckets, client.ProvisionBucket.
			Create().
			SetMinAmount(decimal.NewFromFloat(1000.00)).
			SetMaxAmount(decimal.NewFromFloat(50000.00)).
			SetCurrency(currency),
		)

		_, err := client.ProvisionBucket.
			CreateBulk(sampleBuckets...).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding provision buckets: %w", err)
		}
	}

	// Seed Fiat Currencies
	currencies := []types.SupportedCurrencies{
		{Code: "NGN", Decimals: 2, Name: "Nigerian naira", ShortName: "Naira", Symbol: "₦"},
		{Code: "KES", Decimals: 2, Name: "Kenyan shilling", ShortName: "Swahili", Symbol: "/="},
	}

	listedCurrencies := make([]*ent.FiatCurrencyCreate, 0)
	for _, currency := range currencies {

		_, err := client.FiatCurrency.
			Query().
			Where(
				fiatcurrency.IsEnabledEQ(true),
				fiatcurrency.CodeEQ(currency.Code),
			).
			Only(ctx)
		if ent.IsNotFound(err) {
			fmt.Printf("Seeding currency - %s\n", currency.Code)
			listedCurrencies = append(listedCurrencies, client.FiatCurrency.Create().
				SetCode(currency.Code).
				SetShortName(currency.ShortName).
				SetSymbol(currency.Symbol).
				SetName(currency.Name),
			)
		}
	}

	if len(listedCurrencies) > 0 {
		_, err = client.FiatCurrency.CreateBulk(listedCurrencies...).Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding fiat currencies: %w", err)
		}
	}

	return nil
}
