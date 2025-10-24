package transip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

// data object for DNS entry, used in getting and updating DNS records
type DNSEntry struct {
	Name    string `json:"name"`
	Expire  int    `json:"expire"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type GetDNSEntries struct {
	DNSEntries []DNSEntry `json:"dnsEntries"`
}

func CheckDNSIP(domain string, names []string, ip string, token string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.transip.nl/v6/domains/%s/dns", domain), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var res GetDNSEntries
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return false, err
		}

		for _, dnsEntry := range res.DNSEntries {
			if slices.Contains(names, dnsEntry.Name) && dnsEntry.Content != ip {
				return false, nil
			}
		}

		return true, nil
	}

	return false, fmt.Errorf("got HTTP response with status code %d, expected %d", resp.StatusCode, http.StatusOK)
}

type PatchDNSEntry struct {
	DNSEntry DNSEntry `json:"dnsEntry"`
}

func UpdateIP(domain, name, ip string, token string) error {
	dnsEntry := PatchDNSEntry{
		DNSEntry: DNSEntry{
			Name:    name,
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
