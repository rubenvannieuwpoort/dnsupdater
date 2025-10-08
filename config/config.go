package config

import (
	"fmt"
	"os"
	"strconv"
)

const LoginEnvVar = "LOGIN"
const DnsDomainEnvVar = "DNS_DOMAIN"
const DnsRecordNameEnvVar = "DNS_RECORD"
const CheckIntervalEnvVar = "CHECK_INTERVAL"
const TokenTTLEnvVar = "TOKEN_TTL"

type Config struct {
	Login                string
	DNSDomain            string
	DNSRecordName        string
	CheckIntervalSeconds int
	TokenTTLSeconds      int
}

func Get() Config {
	return Config{
		Login:                getEnv(LoginEnvVar),
		DNSDomain:            getEnv(DnsDomainEnvVar),
		DNSRecordName:        getEnv(DnsRecordNameEnvVar),
		CheckIntervalSeconds: getEnvAsInt(CheckIntervalEnvVar),
		TokenTTLSeconds:      getEnvAsInt(TokenTTLEnvVar),
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
