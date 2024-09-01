package notifications

import "time"

type Transactioin struct {
	ID                 string    `json:"id"`
	CreateDate         time.Time `json:"createDate"`
	UpdateDate         time.Time `json:"updateDate"`
	WalletId           string    `json:"walletId"`
	Amounts            []string  `json:"amounts"`
	NftTokenIds        []string  `json:"nftTokenIds"`
	BlockChain         string    `json:"blockchain"`
	DestinationAddress string    `json:"destinationAddress"`
	State              string    `json:"state"`
	TokenId            string    `json:"tokenId"`
	TransactionType    string    `json:"transactionType"`
	TxHash             string    `json:"txHash"`

	ErrorReason string `json:"errorReason"`
	// ErrorDetails map[string]interface{} `json:"errorDetails"`
}
