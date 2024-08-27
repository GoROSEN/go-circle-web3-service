// devwallets api for developer controlled wallets
package devwallets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CircleDevWalletsService developer controlled wallets(DCW) service
type CircleDevWalletsService struct {
	Host   string
	ApiKey string
}

// NewCircleDevWalletsService create a new instance of DCW service
func NewCircleDevWalletsService(host, apikey string) *CircleDevWalletsService {
	return &CircleDevWalletsService{
		Host:   host,
		ApiKey: apikey,
	}
}

// CreateWalletSet create a wallet set by given keys and walletset name
func (s *CircleDevWalletsService) CreateWalletSet(idempotencyKey, entitySecretCipherText, walletSetName string) (string, error) {

	url := fmt.Sprintf("%v/v1/w3s/developer/walletSets", s.Host)

	payloadObj := struct {
		IdempotencyKey         string `json:"idempotencyKey"`
		EntitySecretCipherText string `json:"entitySecretCipherText"`
		Name                   string `json:"name"`
	}{
		IdempotencyKey:         idempotencyKey,
		EntitySecretCipherText: entitySecretCipherText,
		Name:                   walletSetName,
	}

	return s.httpPost(url, payloadObj)

}

// CreateWallet create a wallet in given blockchains
func (s *CircleDevWalletsService) CreateWallet(idempotencyKey, entitySecretCipherText, walletSetId string, blockChains []string, count int) (string, error) {

	url := fmt.Sprintf("%v/v1/w3s/developer/wallets", s.Host)

	payloadObj := struct {
		IdempotencyKey         string   `json:"idempotencyKey"`
		EntitySecretCipherText string   `json:"entitySecretCipherText"`
		Blockchains            []string `json:"blockchains"`
		Count                  int      `json:"count"`
		WalletSetId            string   `json:"walletSetId"`
	}{
		IdempotencyKey:         idempotencyKey,
		EntitySecretCipherText: entitySecretCipherText,
		Blockchains:            blockChains,
		Count:                  count,
		WalletSetId:            walletSetId,
	}

	return s.httpPost(url, payloadObj)
}

func (s *CircleDevWalletsService) GetWalletBalance(walletId string) (string, error) {

	url := fmt.Sprintf("%v/v1/w3s/wallets/%v/balances", s.Host, walletId)
	return s.httpGet(url, "")
}

func (s *CircleDevWalletsService) SendTransaction(idempotencyKey, entitySecretCipherText, destination, tokenId, walletId string, amounts []string) (string, error) {

	url := fmt.Sprintf("%v/v1/w3s/developer/transactions/transfer", s.Host)

	payloadObj := struct {
		IdempotencyKey         string   `json:"idempotencyKey"`
		EntitySecretCipherText string   `json:"entitySecretCipherText"`
		Amounts                []string `json:"amounts"`
		FeeLevel               string   `json:"feeLevel"`
		TokenId                string   `json:"tokenId"`
		WalletId               string   `json:"walletId"`
		DestinationAddress     string   `json:"destinationAddress"`
	}{
		IdempotencyKey:         idempotencyKey,
		EntitySecretCipherText: entitySecretCipherText,
		Amounts:                amounts,
		FeeLevel:               "HIGH",
		TokenId:                tokenId,
		WalletId:               walletId,
		DestinationAddress:     destination,
	}

	return s.httpPost(url, payloadObj)
}

// httpPost internal method for http post restful calling
func (s *CircleDevWalletsService) httpPost(url string, payloadData interface{}) (string, error) {

	bstr, _ := json.Marshal(payloadData)
	payload := bytes.NewReader(bstr)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}

func (s *CircleDevWalletsService) httpGet(url, queries string) (string, error) {

	var fullUrl string
	if len(queries) > 0 {
		fullUrl = url + "?" + queries
	} else {
		fullUrl = url
	}
	req, _ := http.NewRequest("GET", fullUrl, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}
