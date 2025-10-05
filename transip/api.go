package transip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
)

// data object for DNS entry, used in getting and updating DNS records
type DNSEntry struct {
	Name    string `json:"name"`
	Expire  int    `json:"expire"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

var login, domain string

func init() {
	cnt = rand.Uint64()

	login = os.Getenv("LOGIN")
	if login == "" {
		panic(errors.New("LOGIN environment variable not set"))
	}

	domain = os.Getenv("DOMAIN")
	if domain == "" {
		panic(errors.New("DOMAIN environment variable not set"))
	}

	RefreshToken()
}

type GetDNSEntries struct {
	DNSEntries []DNSEntry `json:"dnsEntries"`
}

func GetDNSIP() (string, error) {
	m.Lock()
	token := _token
	m.Unlock()

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.transip.nl/v6/domains/%s/dns", domain), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var res GetDNSEntries
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return "", err
		}

		for _, dnsEntry := range res.DNSEntries {
			if dnsEntry.Name == "@" {
				return dnsEntry.Content, nil
			}
		}

		return "", errors.New(`no DNS entry for "@"`)
	}

	return "", fmt.Errorf("got HTTP response with status code %d, expected %d", resp.StatusCode, http.StatusOK)
}

type PatchDNSEntry struct {
	DNSEntry DNSEntry `json:"dnsEntry"`
}

func UpdateIP(ip string) error {
	m.Lock()
	token := _token
	m.Unlock()

	dnsEntry := PatchDNSEntry{
		DNSEntry: DNSEntry{
			Name:    "@",
			Expire:  3600,
			Type:    "A",
			Content: ip,
		},
	}

	jsonBytes, err := json.Marshal(dnsEntry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.transip.nl/v6/domains/%s/dns", domain), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
