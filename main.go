package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rubenvannieuwpoort/dnsupdater/config"
	"github.com/rubenvannieuwpoort/dnsupdater/ipify"
	"github.com/rubenvannieuwpoort/dnsupdater/transip"
)

func main() {
	var cfg = config.Get()

	ticker := time.NewTicker(time.Duration(cfg.CheckIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		err := check(cfg)
		if err != nil {
			log.Println(err)
		}

		// wait for the ticker
		_ = <-ticker.C
	}
}

func check(cfg config.Config) error {
	publicIP, err := ipify.GetPublicIP()
	if err != nil {
		return fmt.Errorf("error getting public IP: %v", err)
	}
	log.Printf("got public IP address %s\n", publicIP)

	token, err := transip.GetToken(cfg.Login, cfg.TokenTTLSeconds, cfg.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("error getting token: %v\n", err)
	}
	log.Println("received token for TransIP API")

	ipMatches, err := transip.CheckDNSIP(cfg.DNSDomain, cfg.DNSRecordNames, publicIP, token)
	if err != nil {
		return fmt.Errorf("error getting DNS IP address: %v\n", err)
	}

	if ipMatches {
		log.Printf("IP addresses match, nothing to be done")
	} else {
		log.Printf("IP address needs to be updated")

		for _, name := range cfg.DNSRecordNames {
			err = transip.UpdateIP(cfg.DNSDomain, name, publicIP, token)

			if err != nil {
				return fmt.Errorf("error updating DNS: %v", err)
			}
		}

		log.Print("updated succesfully")
	}

	return nil
}
