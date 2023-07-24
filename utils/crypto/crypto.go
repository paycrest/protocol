package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/paycrest/paycrest-protocol/config"
	"golang.org/x/crypto/bcrypt"
)

var authConf = config.AuthConfig()
var serverConf = config.ServerConfig()

// CheckPasswordHash is a function to compare provided password with the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Encrypt is a function to encrypt plaintext using AES encryption algorithm
func Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(authConf.Secret))
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt is a function to decrypt ciphertext using AES encryption algorithm
func Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(authConf.Secret))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// GenerateAccountFromIndex generates a crypto wallet account from HD wallet mnemonic
func GenerateAccountFromIndex(accountIndex int) (string, string, error) {
	//added code to test generate addrress
	mnemonic := serverConf.HDWalletMnemonic

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", "", fmt.Errorf("failed to create wallet from mnemonic: %w", err)
	}

	path, err := hdwallet.ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", accountIndex))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive account: %w", err)
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to get private key: %w", err)
	}

	privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))
	address := account.Address.Hex()

	return address, privateKeyHex, nil
}
