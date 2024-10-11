package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"
)

func isValidUrl(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func isUrlReachable(testUrl string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(testUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}
	return false
}

func generateShortKey() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	return token
}
