# DNS updater

A personal Go program that solves the problem of pointing a DNS entry to a server on a home network that is behind a NAT while
- Having an ISP (Odido) that only provides dynamic public IP addresses
- Having a DNS provider (TransIP) that does not support dynamic DNS

Luckily, TransIP does provide an API that can be used to automatically update the DNS entry whenever my public IP address changes.

Now, writing a custom program for this is probably about 100 times more work than just setting up dynamic DNS, and there are countless existing scripts and programs that already solve this problem. Naturally, I did what any self-respecting programmer would do: completely ignore all of them and write my own from scratch.

This is tailored for my specific use case and probably isn't directly usable for anyone else (except maybe as a reference). This program is based on the [TransIP API documentation](https://api.transip.nl/rest/docs.html).

The program uses the ipify API to fetch your current public IP address and updates the DNS A record for the specified domain via the TransIP API if it has changed. It runs continuously, checking for IP changes periodically.

Usage:
```
LOGIN=my_transip_username DOMAIN=mydomain.com go run main.go
```

The TransIP private key should be provided in `.secrets/private.pem` (or the `PRIVATE_KEY_PATH` constant in `transip/sign.go` should be updated to point to the private key file).
