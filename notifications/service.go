package notifications

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/martian/log"
)

type NotificationPublicKey struct {
	ID        string    `json:"id"`
	Algorithm string    `json:"algorithm"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createDate"`
}

type NotificationService struct {
	Host   string
	ApiKey string
}

func (ns *NotificationService) GetPublicKey(keyid string) (*NotificationPublicKey, error) {

	var result NotificationPublicKey

	url := fmt.Sprintf("%v/v2/notifications/publicKey/%v", ns.Host, keyid)

	client := resty.New()
	if response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", ns.ApiKey)).
		// SetResult(&result).
		Get(url); err != nil {
		log.Errorf("calling GetNotificationPublicKey, service error: %v", err)
		return nil, err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	return &result, nil
}
