package storage

import (
	"context"
	"flag"
	"fmt"

	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/types"
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

	// Seed Tokens
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

	_, err = client.Token.
		Create().
		SetSymbol("6TEST").
		SetContractAddress("0x3870419Ba2BBf0127060bCB37f69A1b1C090992B").
		SetDecimals(6).
		SetNetwork(network).
		SetIsEnabled(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed seeding token: %w", err)
	}

	// Seed Fiat Currencies and Provision Buckets
	currencies := []types.SupportedCurrencies{
		{Code: "NGN", Decimals: 2, Name: "Nigerian naira", ShortName: "Naira", Symbol: "₦", MarketRate: decimal.NewFromFloat(930.00)},
		{Code: "KES", Decimals: 2, Name: "Kenyan shilling", ShortName: "Swahili", Symbol: "/=", MarketRate: decimal.NewFromFloat(151.45)},
	}
	sampleBuckets := make([]*ent.ProvisionBucketCreate, 0, 6)

	for _, currencyVal := range currencies {
		currency, err := client.FiatCurrency.
			Query().
			Where(
				fiatcurrency.IsEnabledEQ(true),
				fiatcurrency.CodeEQ(currencyVal.Code),
			).
			Only(ctx)
		if ent.IsNotFound(err) {
			currency, _ = client.FiatCurrency.
				Create().
				SetCode(currencyVal.Code).
				SetShortName(currencyVal.ShortName).
				SetSymbol(currencyVal.Symbol).
				SetName(currencyVal.Name).
				SetMarketRate(currencyVal.MarketRate).
				Save(ctx)
		}

		createProvisionBucket := func(min, max float64) *ent.ProvisionBucketCreate {
			return client.ProvisionBucket.
				Create().
				SetMinAmount(decimal.NewFromFloat(min)).
				SetMaxAmount(decimal.NewFromFloat(max)).
				SetCurrency(currency)
		}

		sampleBuckets = append(sampleBuckets, createProvisionBucket(20000001.00, 100000000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(5000001.00, 20000000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(2000001.00, 5000000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(500001.00, 2000000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(50001.00, 500000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(1000.00, 50000.00))
	}

	for _, bucket := range sampleBuckets {
		_, err := bucket.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding provision bucket: %w", err)
		}
	}

	return nil
}
