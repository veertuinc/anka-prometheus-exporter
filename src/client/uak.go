package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UAK struct {
	ID        string
	KeyPath   string
	KeyString string
}

func setUpUAK(uak UAK, controllerAddress string) (string, error) {
	encodedData, err := tap(uak, controllerAddress)
	return encodedData, err
}

func tap(uak UAK, controllerAddress string) (string, error) {

	// Send a POST request to /hand endpoint
	resp, err := http.Post(controllerAddress+"/tap/v1/hand", "application/json", bytes.NewBuffer([]byte(`{"id": "`+uak.ID+`"}`)))
	if err != nil {
		return "", fmt.Errorf("error while sending request to /hand endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response body: %v", err)
	}

	// Base64 decode the response body
	decodedBody, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return "", fmt.Errorf("error while decoding response body (check your id and key path): %v", err)
	}

	// Decrypt the decoded body using the user's private key
	var privateKey *rsa.PrivateKey
	var privateKeyBytes []byte

	if uak.KeyPath != "" {
		privateKeyBytes, err = os.ReadFile(uak.KeyPath)
		if err != nil {
			return "", fmt.Errorf("failed to read private key file: %v", err)
		}
		block, _ := pem.Decode(privateKeyBytes)
		if block == nil {
			return "", errors.New("failed to parse PEM block containing the private key")
		}
		privateKeyBytes = block.Bytes
	} else if uak.KeyString != "" {
		privateKeyBytes, err = base64.StdEncoding.DecodeString(uak.KeyString)
		if err != nil {
			return "", fmt.Errorf("failed decoding private key: %w", err)
		}
	} else {
		return "", fmt.Errorf("no uak path or string specified")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("unable to parse RSA private key: %v", err)
	}
	decryptedBody, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, decodedBody, nil)
	if err != nil {
		return "", fmt.Errorf("rsa decryption error: %v", err)
	}

	// Send a POST request to /shake endpoint
	resp, err = http.Post(controllerAddress+"/tap/v1/shake", "application/json", bytes.NewBuffer([]byte(`{"id": "`+uak.ID+`", "secret": "`+string(decryptedBody)+`"}`)))
	if err != nil {
		return "", fmt.Errorf("error while sending request to /shake endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response body: %v", err)
	}

	// Extract the data field from the response body
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	data := result["data"]

	// Convert the data field to a JSON string
	dataBytes, _ := json.Marshal(data)
	dataString := string(dataBytes)

	// Base64 encode the JSON string
	encodedData := base64.StdEncoding.EncodeToString([]byte(dataString))

	// // Set the Authorization header for subsequent requests
	// client := &http.Client{}
	// req, err := http.NewRequest("GET", controllerAddress+"/api/v1/status", nil)
	// if err != nil {
	// 	return "", fmt.Errorf("error while creating new request: %v", err)
	// }
	// req.Header.Set("Authorization", "Bearer "+encodedData)

	// // Send the request
	// resp, err = client.Do(req)
	// if err != nil {
	// 	return "", fmt.Errorf("error while sending request: %v", err)
	// }
	// defer resp.Body.Close()

	// // Read the response body
	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("error while reading response body: %v", err)
	// }

	// // Unmarshal the response body
	// var responseBody map[string]interface{}
	// json.Unmarshal(body, &responseBody)

	// log.Info(fmt.Sprintf("[auth::uak] Response from /status endpoint: %s", string(body)))

	return encodedData, nil
}
