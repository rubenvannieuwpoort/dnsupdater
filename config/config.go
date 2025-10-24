package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const LoginEnvVar = "LOGIN"
const DnsDomainEnvVar = "DNS_DOMAIN"
const DnsRecordNameEnvVar = "DNS_RECORDS"
const CheckIntervalEnvVar = "CHECK_INTERVAL"
const TokenTTLEnvVar = "TOKEN_TTL"
const PrivateKeyPathEnvVar = "PRIVATE_KEY_PATH"

type Config struct {
	Login                string
	DNSDomain            string
	DNSRecordNames       []string
	CheckIntervalSeconds int
	TokenTTLSeconds      int
	PrivateKeyPath       string
}

func Get() Config {
	dnsRecordStr := getEnv(DnsRecordNameEnvVar)
	return Config{
		Login:                getEnv(LoginEnvVar),
		DNSDomain:            getEnv(DnsDomainEnvVar),
		DNSRecordNames:       strings.Split(dnsRecordStr, ","),
		CheckIntervalSeconds: getEnvAsInt(CheckIntervalEnvVar),
		TokenTTLSeconds:      getEnvAsInt(TokenTTLEnvVar),
		PrivateKeyPath:       getEnv(PrivateKeyPathEnvVar),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Errorf("required environment variable %s is not set", key))
	}
	return value
}

func getEnvAsInt(key string) int {
	valueStr := getEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Errorf("invalid integer value for %s: %v", key, err))
	}
	return value
}
