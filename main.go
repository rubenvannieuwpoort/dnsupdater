package main

import (
	"log"
	"time"

	"github.com/rubenvannieuwpoort/dnsupdater/ipify"
	"github.com/rubenvannieuwpoort/dnsupdater/transip"
)

func main() {
	// refresh API token in a separate goroutine
	go periodicallyRefreshToken(9 * time.Minute)

	lastIP, err := transip.GetDNSIP()
	if err != nil {
		log.Printf("error getting IP in DNS: %v", err)
	}

	for {
		publicIP, err := ipify.GetPublicIP()
		if err != nil {
			log.Printf("error getting public IP: %v", err)
			continue
		}

		log.Printf("got public IP %s", publicIP)

		if publicIP != lastIP {
			log.Printf("public IP \"%s\" does not match last known IP \"%s\" in DNS, updating...", publicIP, lastIP)
			err = transip.UpdateIP(publicIP)

			if err != nil {
				log.Printf("error updating DNS: %v", err)
			} else {
				log.Print("updated succesfully")
			}
			lastIP = publicIP
		} else {
			log.Print("public IP matches last known IP in DNS, no action needed")
		}

		time.Sleep(1 * time.Minute)
	}
}

func periodicallyRefreshToken(refreshInterval time.Duration) {
	for {
		time.Sleep(refreshInterval)

		// update API token
		err := transip.RefreshToken()
		if err != nil {
			log.Printf("error refreshing token: %v", err)
		} else {
			log.Printf("updated token succesfully")
		}
	}
}
