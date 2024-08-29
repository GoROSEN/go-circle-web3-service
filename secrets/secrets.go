package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/google/martian/log"
)

func GenerateEntitySecret() []byte {
	mainBuff := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, mainBuff)
	if err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return mainBuff
}

func GetPublickKey(host, apikey string) (string, error) {

	var result struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"message,omitempty"`
		Data struct {
			PublicKey string `json:"publicKey"`
		} `json:"data,omitempty"`
	}

	url := fmt.Sprintf("%v/v1/w3s/config/entity/publicKey", host)

	client := resty.New()
	if response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", apikey)).
		// SetResult(&result).
		Get(url); err != nil {
		log.Errorf("calling get public service error: %v", err)
		return "", err
	} else {
		// resty doesn't unmarshal the response to result, i don't know why
		json.Unmarshal(response.Body(), &result)
	}

	if result.Code != 0 {
		log.Errorf("get public service got error code: %v, reason: %v", result.Code, result.Msg)
		return "", fmt.Errorf(result.Msg)
	}

	return result.Data.PublicKey, nil
}

func EncryptEntitySecret(hexEncodedEntitySecret, rsaPublicKeyString string) (string, error) {
	entitySecret, err := hex.DecodeString(hexEncodedEntitySecret)
	if err != nil {
		return "", err
	}

	if len(entitySecret) != 32 {
		return "", errors.New("invalid entity secret")
	}

	pubKey, err := parseRsaPublicKeyFromPem([]byte(rsaPublicKeyString))
	if err != nil {
		return "", err
	}

	cipher, err := encryptOAEP(pubKey, entitySecret)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

// ParseRsaPublicKeyFromPem parse rsa public key from pem.
func parseRsaPublicKeyFromPem(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
	}
	return nil, errors.New("key type is not rsa")
}

// EncryptOAEP rsa encrypt oaep.
func encryptOAEP(pubKey *rsa.PublicKey, message []byte) (ciphertext []byte, err error) {
	random := rand.Reader
	ciphertext, err = rsa.EncryptOAEP(sha256.New(), random, pubKey, message, nil)
	if err != nil {
		return nil, err
	}
	return
}
