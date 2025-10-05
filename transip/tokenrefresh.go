package transip

// this file provides the functionality

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"sync"
)

const TOKEN_EXPIRATION_TIME = "30 minutes"

// the _token can potentially be concurrently accessed so (theoretically) we need a mutex
// to prevent data races (I'm sure in practice it'll work fine without the mutex)
var m sync.Mutex
var _token string

// used as a suffix in the API token name to prevent duplicates
var cnt uint64

type UpdateTokenRequest struct {
	Login          string `json:"login"`
	Nonce          string `json:"nonce"`
	ReadOnly       bool   `json:"read_only"`
	ExpirationTime string `json:"expiration_time"`
	Label          string `json:"label"`
	GlobalKey      bool   `json:"global_key"`
}

type UpdateTokenResponse struct {
	Token string `json:"token"`
}

func RefreshToken() error {
	requestBody := UpdateTokenRequest{
		Login:          login,
		ReadOnly:       false,
		ExpirationTime: TOKEN_EXPIRATION_TIME,
		GlobalKey:      true,
	}

	requestBody.Nonce = fmt.Sprintf("%d", rand.IntN(math.MaxInt))
	requestBody.Label = fmt.Sprintf("dnsupdater%d", cnt)
	cnt++

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest("POST", "https://api.transip.nl/v6/auth", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	signature, err := sign(body)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", signature)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var res UpdateTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return err
		}

		m.Lock()
		_token = res.Token
		m.Unlock()

		return nil
	}

	return fmt.Errorf("got HTTP response with status code %d, expected %d", resp.StatusCode, http.StatusCreated)
}
