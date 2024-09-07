# go-circle-web3-service
a Golang sdk for circle web3 service

Currently we have only implemented required APIs for developer controlled wallets. You can check the [example](examples/basic.go) for a completed usages.

To use it in your project, you can simply do this on your golang project:

```
go get github.com/GoROSEN/go-circle-web3-service
```

then import and use it on your codes.

```

import (
  "github.com/GoROSEN/go-circle-web3-service"
)

func foobar() {

	host := "https://api.circle.com" // endpoint of circle web3 service
	apikey := "<your api key of circle web3 service here>"
	secret := "<your secret hex string which has been registered on the console of circle web3 service here>"
	pubkey := "<your public key from the the console of circle web3 service>"

  service := devwallets.NewCircleDevWalletsService(host, apikey, pubkey, secret)

  service.CreateWalletSet(...)
  service.CreateWallet(...)
  service.GetWalletBalanceSimple(wallet-id)
  service.GetWalletBalance(...)
  service.SendTransaction(...)
}

```