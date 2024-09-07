package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoROSEN/go-circle-web3-service/devwallets"
)

// usage:
// go run examples/basic.go --host https://api.circle.com --apikey <your_api_key_here> --secret <your_registered_secret_here> --pubkey <your_public_key_pem_file>
func main() {

	host := flag.String("host", "https://api.circle.com", "endpoint of circle web3 service")
	apikey := flag.String("apikey", "", "your api key of circle web3 service, starts with 'TEST_API_KEY:' or 'LIVE_API_KEY:'")
	secret := flag.String("secret", "", "your secret hex string which has been registered on the console of circle web3 service. it is a long hex digits looks like 'a2809c388a7f9c1220e19dd7f1ba15c95447a30cfc7c946baf6d36a4b4be087d'")
	pubkeyFile := flag.String("pubkey", "key.pem", "your public key from the the console of circle web3 service. you can both download it from the console, or call GetPublickKey function from the secrets package of this SDK.")

	if apikey == nil || len(*apikey) == 0 {
		fmt.Println("apikey should not be empty")
		return
	}
	if secret == nil || len(*secret) == 0 {
		fmt.Println("secret should not be empty")
		return
	}

	pubkey, err := os.ReadFile(*pubkeyFile)
	if err != nil {
		fmt.Printf("cannot read public key from %v\n", *pubkeyFile)
		return
	}

	// create devwallet service from given parameters
	service := devwallets.NewCircleDevWalletsService(*host, *apikey, string(pubkey), *secret)
	if service == nil {
		fmt.Println("cannot create service")
		return
	}

	// create wallet set
	walletSet, err := service.CreateWalletSet("Basic Set")
	if err != nil {
		fmt.Printf("cannot create wallet set: %v\n", err)
		return
	}

	// create wallets, for test api key, it can only create wallets on devnet or testnet.
	// you have to use live api key to create any wallet on the mainnet
	walletCountToCreate := 1
	wallets, err := service.CreateWallet(walletSet.ID, "EOA", []string{"SOL-DEVNET"}, walletCountToCreate)
	if err != nil {
		fmt.Printf("cannot create wallets: %v\n", err)
		return
	}

	if len(wallets) == 0 {
		fmt.Println("no wallet was created")
		return
	}

	// get every token balances for the wallet
	balances, err := service.GetWalletBalanceSimple(wallets[0].ID)
	if err != nil {
		fmt.Printf("cannot create wallets: %v\n", err)
		return
	}
	for i := range balances {
		fmt.Printf("balance of %v in %v: %v\n", balances[i].Token.Name, wallets[0].Address, balances[i].Amount)
	}

	// send 1 USDC to given address
	tokenId := "8fb3cadb-0ef4-573d-8fcd-e194f961c728" // USDC on SOL-DEVNET
	txhash, err := service.SendTransaction("<set your destination address here>", tokenId, wallets[0].ID, []string{"1.0"})
	if err != nil {
		fmt.Printf("cannot send transaction: %v\n", err)
		return
	}
	fmt.Printf("tx %v has been sent\n", txhash)
}
