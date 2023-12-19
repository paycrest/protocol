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

	// Delete existing data
	_ = client.Network.Delete().ExecX(ctx)
	_ = client.Token.Delete().ExecX(ctx)
	_ = client.FiatCurrency.Delete().ExecX(ctx)
	_ = client.ProvisionBucket.Delete().ExecX(ctx)
	_ = client.User.Delete().ExecX(ctx)
	_ = client.ProviderProfile.Delete().ExecX(ctx)

	// Seed Network
	fmt.Println("seeding network...")
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
	fmt.Println("seeding tokens...")
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
	fmt.Println("fiat currencies and provision buckets...")
	currencies := []types.SupportedCurrencies{
		{Code: "NGN", Decimals: 2, Name: "Nigerian Naira", ShortName: "Naira", Symbol: "₦", MarketRate: decimal.NewFromFloat(930.00)},
		{Code: "KES", Decimals: 2, Name: "Kenyan Shilling", ShortName: "Swahili", Symbol: "KSh", MarketRate: decimal.NewFromFloat(151.45)},
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

	for i, bucketCreate := range sampleBuckets {
		bucket, err := bucketCreate.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding provision bucket: %w", err)
		}

		// Seed users
		fmt.Println("seed users...")
		users := make([]*ent.User, 0, len(sampleBuckets)*3)
		for i := 0; i < len(sampleBuckets)*3; i++ {
			user, err := client.User.
				Create().
				SetFirstName(fmt.Sprintf("User_%d", i)).
				SetLastName("Doe").
				SetEmail(fmt.Sprintf("user_%d@example.com", i)).
				SetPassword("password").
				SetScope("provider sender").
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed creating user: %w", err)
			}
			users = append(users, user)
		}

		fmt.Println("seed provider profiles...")
		currency := bucket.QueryCurrency().FirstX(ctx)
		for j := 0; j < 3; j++ {
			fmt.Println(currency.Code)
			_, err := client.ProviderProfile.
				Create().
				SetTradingName(fmt.Sprintf("Provider_%d", i*3+j)).
				SetUser(users[i*3+j]).
				SetCurrency(currency).
				AddProvisionBuckets(bucket).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed creating provider: %w", err)
			}

			fmt.Println(currency.Code)
		}
	}

	return nil
}
