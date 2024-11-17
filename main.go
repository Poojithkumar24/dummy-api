package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Configuration for API
const (
	UATURL              = "https://apiuat.iifl.in/PayinGateway/Partner/DREMakerV3"
	SubscriptionKey     = "06573965fb254f63bcc4a310b43650fe"
	AppName             = "DREPartner"
	AppVersion          = "1.0"
	EncryptionKey       = "E7D6A56889FEED8F" // Replace with the actual key
	EncryptionKeyLength = 32                 // 256-bit AES key
	EncryptionBlockSize = 16                 // AES block size
)

// RequestHeader defines the structure of the request header
type RequestHeader struct {
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
}

// RequestBody defines the structure of the request body
type RequestBody struct {
	Data string `json:"data"`
}

// APIRequest combines the header and body
type APIRequest struct {
	Head RequestHeader `json:"head"`
	Body RequestBody   `json:"body"`
}

// PlainRequest represents the unencrypted request payload
type PlainRequest struct {
	Location            string  `json:"Location"`
	ProspectNo          string  `json:"ProspectNo"`
	ProductCode         string  `json:"ProductCode"`
	AdvancedEMI         float64 `json:"AdvancedEMI"`
	EMIOTC              float64 `json:"EMI_OTC"`
	BounceChequeCharges float64 `json:"BounceChequeCharges"`
	PenalCharges        float64 `json:"PenalCharges"`
	ProcessingFeeIMD    float64 `json:"ProcessingFee_IMD"`
	PreEMI              float64 `json:"PreEMI"`
	SwapCharges         float64 `json:"SwapCharges"`
	ForeClosure         float64 `json:"ForeClosure"`
	OtherCharges        float64 `json:"OtherCharges"`
	TotalAmount         float64 `json:"TotalAmount"`
	PaymentMode         string  `json:"PaymentMode"`
	BankName            string  `json:"BankName"`
	ChequeNo            string  `json:"ChequeNo"`
	Remarks             string  `json:"Remarks"`
	Reason              string  `json:"Reason"`
	Source              string  `json:"Source"`
	Status              string  `json:"Status"`
}

// EncryptAES256 encrypts data using AES-256 CBC mode
func EncryptAES256(data, key string) (string, error) {
	if len(key) != EncryptionKeyLength {
		return "", fmt.Errorf("invalid encryption key length")
	}

	// Convert the key and plaintext to byte slices
	keyBytes := []byte(key)
	plaintext := []byte(data)

	// Pad plaintext to be a multiple of AES block size
	padding := EncryptionBlockSize - (len(plaintext) % EncryptionBlockSize)
	paddedPlaintext := append(plaintext, bytes.Repeat([]byte{0}, padding)...)

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Generate an initialization vector (IV)
	iv := make([]byte, block.BlockSize()) // Set to zeros (can be random for better security)

	// Perform encryption
	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Encode to Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// CallAPI sends the request to the API and returns the response
func CallAPI(encryptedData string) (*http.Response, error) {
	// Create the request payload
	payload := APIRequest{
		Head: RequestHeader{
			AppName:    AppName,
			AppVersion: AppVersion,
		},
		Body: RequestBody{
			Data: encryptedData,
		},
	}

	// Marshal payload to JSON
	requestBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", UATURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", SubscriptionKey)

	// Perform the HTTP request
	client := &http.Client{}
	return client.Do(req)
}

func main() {
	// Define the API URL
	apiURL := "https://apiuat.iifl.in/PayinGateway/Partner/DREMakerV3"

	// Define the headers
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": "06573965fb254f63bcc4a310b43650fe",
		"ClientId":                  "Paytm",
		"AppName":                   "DREPartner",
		"Content-Type":              "application/json",
	}

	// Define the request payload
	requestBody := `{
		"head": {
			"appName": "DREPartner",
			"appVersion": "1.0"
		},
		"body": {
			"data": "/khw6WiOrS0av6YU1wrtogfnq608dP8ZHOeM95vbi9bgT9OWKfY57b9RbwAIScHHzaBXgRmInw433KSRiIlO0PkMpToJ6rPUxQDQbPGruT4dKfaVxSFaK6zHugm0FG16FJrFJz1obWpc4jfubH2RDZuQfGzF4hvy7bKB9lTxbu1uB7z9S0EbhATb2hA85kc/gehQujqx9x3J6mGXotOdGEnYBmQ6ONEHVeFTIJmhpFrNbTRSxz4ZLLq86b33C1TH+uZjjydwkS5ZtrF+jNrwubqcZZRZvz7iL/sfdPyVv3SImZHR7WnbtK0D2BBhZq73U5d7X0fxwDc/PlGzz8N4BDFxiEeBCRRCl72I+uLxqJMQYS/5fDf54sKJfBR7M2i/YXVnB4Ef0kYHVQ6arvpKKYeKjWkWIUjxpfX+RdpAzUbZNDwVpWizsa/pLrRksQxMv1BU/NFDfYB97OAWKEu2lWfts1gReKwna/VMIMGYVoIFrLpGBSCAX0t0vIMIC0pWSod0uDBvn1FGeqbD3iVvuw7zvp9ET/u0KIXSpZ73cG0/21WTVhKxGkDK4S7PZ2W3dEbFJjUDwaptk+oGtjKFD2fzoyUGwC0dVMVKlyFaRbh36thn4icaD8+fK7VBAWhFueNmXVtjrnYpgS3pOm7rkGUnaTZ/bDbWVQxBCTiMvnvbIj2ocSn8wg82yhvltET+iG1f0PxEe3cQZ7XUYRibG78RjZcKelpehIgw786x+TC7mDoPwJnlD7PVex793Pqn04oZdxCMS4v0rR5buYA17bSTp0/xOvMF9viZM5sIpF4="
		}
	}`

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}
