package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateReceive(t *testing.T) {
	// Mock the server config
	serverConf.HDWalletMnemonic = "media nerve fog identify typical physical aspect doll bar fossil frost because"

	// set the expected account index and address
	expectedAccountIndex := 1
	expectedAddress := "0xc60F0aDe1483fa6A355f32E0d3406127C49d4d7f"

	// Call the GenerateReceiveAddress Function
	address,privateKey, err := GenerateReceiveAddress(expectedAccountIndex)

	// error checker 
	assert.NoError(t, err, "unexpected error")

	// Assert the generated address
	assert.Equal(t, expectedAddress, address, "incorrect address")
	assert.NotEmpty(t, privateKey, "private key should not be empty")
}
