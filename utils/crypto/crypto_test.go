package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEOA(t *testing.T) {
	// Mock the server config
	cryptoConf.HDWalletMnemonic = "media nerve fog identify typical physical aspect doll bar fossil frost because"

	t.Run("evm account creation", func(t *testing.T) {
		// Set the expected account index and address
		expectedAccountIndex := 1
		expectedAddress := "0xc60F0aDe1483fa6A355f32E0d3406127C49d4d7f"

		// Call the GenerateAccountFromIndex Function
		address, privateKey, err := GenerateAccountFromIndex(expectedAccountIndex)
		assert.NoError(t, err, "unexpected error")

		// Assert the generated address
		assert.Equal(t, expectedAddress, address.Hex(), "incorrect address")
		assert.NotEmpty(t, privateKey, "private key should not be empty")
	})

	t.Run("tron account creation", func(t *testing.T) {
		// Set the expected account index and address
		expectedAccountIndex := 1
		expectedAddress := "TFR3TTx4YzWwNoqmcuVEi477PJoSyo9zwx"

		// Call the GenerateTronAccountFromIndex Function
		wallet, err := GenerateTronAccountFromIndex(expectedAccountIndex)
		assert.NoError(t, err, "unexpected error")

		// Assert the generated address
		assert.Equal(t, expectedAddress, wallet.AddressBase58, "incorrect address")
		assert.NotEmpty(t, wallet.PrivateKey, "private key should not be empty")
	})
}
