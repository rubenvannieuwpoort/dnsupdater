package transip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
)

type UpdateTokenRequest struct {
	Label          string `json:"label"`
	Login          string `json:"login"`
	Nonce          string `json:"nonce"`
	ReadOnly       bool   `json:"read_only"`
	ExpirationTime string `json:"expiration_time"`
	GlobalKey      bool   `json:"global_key"`
}

type UpdateTokenResponse struct {
	Token string `json:"token"`
}

func GetToken(login string, ttl int, privateKeyPath string) (string, error) {
	requestBody := UpdateTokenRequest{
		Label:          "dnsupdater",
		Login:          login,
		Nonce:          fmt.Sprintf("%d", rand.IntN(math.MaxInt)),
		ReadOnly:       false,
		ExpirationTime: fmt.Sprintf("%d seconds", ttl),
		GlobalKey:      true,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling body from response from api.transip.nl to JSON: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "https://api.transip.nl/v6/auth", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %v", err)
	}

	signature, err := sign(body, privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error signing body: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", signature)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending POST request to api.transip.nl: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("got HTTP response with status code %d POST request to api.transip.nl, expected %d", resp.StatusCode, http.StatusCreated)
	}

	var res UpdateTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON from response from api.transip.nl: %v", err)
	}

	return res.Token, nil
}
