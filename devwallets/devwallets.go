// devwallets api for developer controlled wallets
package devwallets

import (
	"encoding/json"
	"fmt"

	"github.com/GoROSEN/go-circle-web3-service/secrets"
	"github.com/go-resty/resty/v2"
	"github.com/gofrs/uuid"
	"github.com/google/martian/log"
)

// CircleDevWalletsService developer controlled wallets(DCW) service
type CircleDevWalletsService struct {
	Host      string
	ApiKey    string
	PublicKey string
	Secret    string
}

// NewCircleDevWalletsService create a new instance of DCW service
func NewCircleDevWalletsService(host, apikey, pubkey, secret string) *CircleDevWalletsService {
	return &CircleDevWalletsService{
		Host:      host,
		ApiKey:    apikey,
		PublicKey: pubkey,
		Secret:    secret,
	}
}

// CreateWalletSet create a wallet set by given keys and walletset name
func (s *CircleDevWalletsService) CreateWalletSet(walletSetName string) (*DevWalletSet, error) {

	entitySecretCipherText, _ := secrets.EncryptEntitySecret(s.Secret, s.PublicKey)
	tmp, _ := uuid.NewV4()
	idempotencyKey := tmp.String()

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

	var result struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"message,omitempty"`
		Data struct {
			WalletSet DevWalletSet `json:"walletSet"`
		} `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}

	client := resty.New()
	if response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey)).
		SetBody(&payloadObj).
		// SetResult(&result).
		Post(url); err != nil {
		log.Errorf("calling CreateWalletSet, service error: %v", err)
		return nil, err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	if result.Code != 0 {
		log.Errorf("CreateWalletSet got error code: %v, reason: %v", result.Code, result.Msg)
		return nil, fmt.Errorf(result.Msg)
	}

	return &result.Data.WalletSet, nil
}

// CreateWallet create a wallet in given blockchains
func (s *CircleDevWalletsService) CreateWallet(walletSetId, accountType string, blockChains []string, count int) ([]Wallet, error) {

	entitySecretCipherText, _ := secrets.EncryptEntitySecret(s.Secret, s.PublicKey)
	tmp, _ := uuid.NewV4()
	idempotencyKey := tmp.String()

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

	var result struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"message,omitempty"`
		Data struct {
			Wallets []Wallet `json:"wallets"`
		} `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}

	client := resty.New()
	if response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey)).
		SetBody(&payloadObj).
		// SetResult(&result).
		Post(url); err != nil {
		log.Errorf("calling CreateWallet, service error: %v", err)
		return nil, err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	if result.Code != 0 {
		log.Errorf("CreateWallets got error code: %v, reason: %v", result.Code, result.Msg)
		return nil, fmt.Errorf(result.Msg)
	}

	return result.Data.Wallets, nil
}

// GetWalletBalanceSimple get wallet tokens balance
func (s *CircleDevWalletsService) GetWalletBalanceSimple(walletId string) ([]TokenBalance, error) {
	return s.GetWalletBalance(walletId, "", "", "", "", "", false, 0)
}

// GetWalletBalance get paged tokens balances
func (s *CircleDevWalletsService) GetWalletBalance(walletId, name, tokenAddress, standard, pageBefore, pageAfter string, includeAll bool, pageSize int) ([]TokenBalance, error) {

	url := fmt.Sprintf("%v/v1/w3s/wallets/%v/balances", s.Host, walletId)

	var result struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"message,omitempty"`
		Data struct {
			TokenBalance []TokenBalance `json:"tokenBalances,omitempty"`
		} `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}

	client := resty.New()
	r := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey))
	if includeAll {
		r.SetQueryParam("includeAll", "true")
	}
	if len(name) > 0 {
		r.SetQueryParam("name", name)
	}
	if len(tokenAddress) > 0 {
		r.SetQueryParam("tokenAddress", tokenAddress)
	}
	if len(standard) > 0 {
		r.SetQueryParam("standard", standard)
	}
	if len(pageBefore) > 0 {
		r.SetQueryParam("pageBefore", pageBefore)
	} else if len(pageAfter) > 0 {
		r.SetQueryParam("pageAfter", pageAfter)
	}
	if pageSize > 0 && pageSize <= 50 {
		r.SetQueryParam("pageSize", fmt.Sprintf("%v", pageSize))
	}
	if response, err := r.
		// SetResult(&result).
		Get(url); err != nil {
		log.Errorf("calling get public service error: %v", err)
		return nil, err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	return result.Data.TokenBalance, nil
}

// SendTransaction send a transaction to circle web 3 service
func (s *CircleDevWalletsService) SendTransaction(destination, tokenId, walletId string, amounts []string) (string, error) {

	entitySecretCipherText, _ := secrets.EncryptEntitySecret(s.Secret, s.PublicKey)
	tmp, _ := uuid.NewV4()
	idempotencyKey := tmp.String()

	var result struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"message,omitempty"`
		Data struct {
			Id    string `json:"id"`
			State string `json:"state"` // INITIATED PENDING_RISK_SCREENING DENIED QUEUED SENT CONFIRMED COMPLETE FAILED CANCELLED
		} `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}

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
		FeeLevel:               "LOW",
		TokenId:                tokenId,
		WalletId:               walletId,
		DestinationAddress:     destination,
	}

	client := resty.New()
	if response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", s.ApiKey)).
		SetBody(&payloadObj).
		// SetResult(&result).
		Post(url); err != nil {
		log.Errorf("calling CreateWallet, service error: %v", err)
		return "", err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	if result.Code != 0 {
		log.Errorf("CreateWallets got error code: %v, reason: %v", result.Code, result.Msg)
		return "", fmt.Errorf(result.Msg)
	}

	return result.Data.Id, nil
}
