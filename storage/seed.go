package storage

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"time"

	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/token"
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
	_ = client.ProviderOrderToken.Delete().ExecX(ctx)
	_ = client.SenderProfile.Delete().ExecX(ctx)

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
		// {Code: "KES", Decimals: 2, Name: "Kenyan Shilling", ShortName: "Swahili", Symbol: "KSh", MarketRate: decimal.NewFromFloat(151.45)},
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
				SetIsEnabled(true).
				Save(ctx)
		}

		createProvisionBucket := func(min, max float64) *ent.ProvisionBucketCreate {
			return client.ProvisionBucket.
				Create().
				SetMinAmount(decimal.NewFromFloat(min)).
				SetMaxAmount(decimal.NewFromFloat(max)).
				SetCurrency(currency)
		}

		sampleBuckets = append(sampleBuckets, createProvisionBucket(5001.00, 50000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(1001.00, 5000.00))
		sampleBuckets = append(sampleBuckets, createProvisionBucket(0.00, 1000.00))
	}

	// Seed users and provider profiles
	fmt.Println("seed users, provider profiles, and order tokens...")

	for i, sampleBucket := range sampleBuckets {
		bucket, err := sampleBucket.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding provision bucket: %w", err)
		}

		for j := 0; j < 2; j++ {
			user, err := client.User.
				Create().
				SetFirstName(fmt.Sprintf("User_%d%d", i, j)).
				SetLastName("Doe").
				SetEmail(fmt.Sprintf("user_%d%d@example.com", i, j)).
				SetPassword("password").
				SetScope("provider sender").
				SetIsEmailVerified(true).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed creating user: %w", err)
			}

			if j == 0 {
				_, err = client.SenderProfile.
					Create().
					SetUser(user).
					SetWebhookURL("https://example.com/webhook").
					SetFeePerTokenUnit(decimal.NewFromFloat(10)).
					SetFeeAddress("0x409689E3008d43a9eb439e7B275749D4a71D8E2D").
					SetRefundAddress("0x409689E3008d43a9eb439e7B275749D4a71D8E2D").
					SetDomainWhitelist([]string{"https://example.com"}).
					SetIsActive(true).
					Save(ctx)
				if err != nil {
					return fmt.Errorf("failed creating sender profile: %w", err)
				}
			}

			currency := bucket.QueryCurrency().OnlyX(ctx)

			provider, err := client.ProviderProfile.
				Create().
				SetTradingName(fmt.Sprintf("Provider_%d%d", i, j)).
				SetHostIdentifier("http://localhost:8001").
				SetUser(user).
				SetIsActive(true).
				SetIsAvailable(true).
				SetCurrencyID(currency.ID).
				SetAddress("123 Main St").
				SetMobileNumber("+2348063000000").
				SetDateOfBirth(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)).
				SetBusinessName("ABC Corporation").
				SetIdentityDocumentType("passport").
				SetIdentityDocument("https://example.com/identity_document.jpg").
				SetBusinessDocument("https://example.com/business_document.pdf").
				AddProvisionBuckets(bucket).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed creating provider: %w", err)
			}

			// Create API Key
			secretKey, err := token.GeneratePrivateKey()
			if err != nil {
				return fmt.Errorf("failed to generate API key: %w", err)
			}
			fmt.Printf("Provider Name: %s, Provider ID: %s, API Secret Key: %s\n", user.FirstName+" "+user.LastName, provider.ID, secretKey)
			encryptedSecret, _ := crypto.EncryptPlain([]byte(secretKey))
			encodedSecret := base64.StdEncoding.EncodeToString(encryptedSecret)

			_, err = client.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetProviderProfile(provider).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create API key: %w", err)
			}

			// Configure tokens
			addresses := []struct {
				Address string `json:"address"`
				Network string `json:"network"`
			}{
				{Address: "0x409689E3008d43a9eb439e7B275749D4a71D8E2D", Network: "polygon-mumbai"},
			}

			_, err = client.ProviderOrderToken.
				Create().
				SetSymbol("6TEST").
				SetConversionRateType("fixed").
				SetFixedConversionRate(decimal.NewFromFloat(1100)).
				SetFloatingConversionRate(decimal.NewFromFloat(0.0)).
				SetMinOrderAmount(bucket.MinAmount).
				SetMaxOrderAmount(bucket.MaxAmount).
				SetAddresses(addresses).
				SetProviderID(provider.ID).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to configure order tokens: %w", err)
			}
		}
	}

	return nil
}
