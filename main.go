package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosmos/go-bip39"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/routers"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

func main() {
	// Set timezone
	conf := config.ServerConfig()
	loc, _ := time.LoadLocation(conf.Timezone)
	time.Local = loc

	// Connect to the database
	DSN := config.DBConfiguration()

	if err := database.DBConnection(DSN); err != nil {
		logger.Fatalf("database DBConnection error: %s", err)
	}

	defer database.GetClient().Close()

	//added code to test generate addrress
	mnemonic := os.Getenv("MNEMONIC")
	if mnemonic == "" {
		log.Fatal("Mnemonic phrase not provided")
	}

	valid := bip39.IsMnemonicValid(mnemonic)
	if !valid {
		log.Fatal("Invalid mnemonic phrase")
	}

	initialIndex := 0

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("Terminating the program...")
		os.Exit(0)
	}()
	// Initialize the services
	receiveAddressService := services.NewReceiveAddressService(wallet, initialIndex)

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :%v", appServer)

	logger.Fatalf("%v", router.Run(appServer))

}
