package devwallets

import "time"

type DevWalletSet struct {
	ID          string    `json:"id"`
	CreateDate  time.Time `json:"createDate"`
	UpdateDate  time.Time `json:"updateDate"`
	CustodyType string    `json:"custodyType"` // DEVELOPER ENDUSER
}

type Wallet struct {
	ID          string    `json:"id"`
	Address     string    `json:"address"`
	Blockchain  string    `json:"blockchain"` // ETH ETH-SEPOLIA AVAX AVAX-FUJI MATIC MATIC-AMOY SOL SOL-DEVNET
	CreateDate  time.Time `json:"createDate"`
	UpdateDate  time.Time `json:"updateDate"`
	CustodyType string    `json:"custodyType"` // DEVELOPER ENDUSER
	Name        string    `json:"name,omitempty"`
	RefId       string    `json:"refId,omitempty"`
	State       string    `json:"state"` // LIVE FROZEN
	UserId      string    `json:"userId,omitempty"`
	WalletSetId string    `json:"walletSetId"`
	AccountType string    `json:"accountType"`       // SCA EOA
	ScaCore     string    `json:"scaCore,omitempty"` // required for SCA wallet
}

type Token struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Standard     string    `json:"standard,omitempty"`
	Blockchain   string    `json:"blockchain"`
	Decimals     int       `json:"decimals,omitempty"`
	Symbol       string    `json:"symbol,omitempty"`
	TokenAddress string    `json:"tokenAddress,omitempty"`
	CreateDate   time.Time `json:"createDate"`
	UpdateDate   time.Time `json:"updateDate"`
}

type TokenBalance struct {
	Amount     string    `json:"amount"`
	Token      Token     `json:"token"`
	UpdateDate time.Time `json:"updateDate"`
}
