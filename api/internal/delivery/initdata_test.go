package delivery

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"testing"
)

func buildInitData(values map[string]string, botToken string) string {
	// build query string without hash
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+url.QueryEscape(values[k]))
	}
	// ParseQuery expects encoded string; we'll create raw string for VerifyInitData
	// But compute data_check_string using unescaped values per algorithm
	partsForCheck := make([]string, 0, len(keys))
	for _, k := range keys {
		partsForCheck = append(partsForCheck, k+"="+values[k])
	}
	dataCheckString := strings.Join(partsForCheck, "\n")
	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))

	// build final query string (URL-encoded values and hash param)
	q := url.Values{}
	for _, k := range keys {
		q.Set(k, values[k])
	}
	q.Set("hash", expected)
	return q.Encode()
}

func TestVerifyInitData_Valid(t *testing.T) {
	botToken := "test_bot_token"
	vals := map[string]string{"id": "12345", "username": "john", "first_name": "John"}
	initData := buildInitData(vals, botToken)
	parsed, err := VerifyInitData(initData, botToken)
	if err != nil {
		t.Fatalf("expected valid initData, got err: %v", err)
	}
	if parsed["id"] != "12345" || parsed["username"] != "john" {
		t.Fatalf("parsed values mismatch: %v", parsed)
	}
}

func TestVerifyInitData_Tampered(t *testing.T) {
	botToken := "test_bot_token"
	vals := map[string]string{"id": "12345", "username": "john"}
	initData := buildInitData(vals, botToken)
	// tamper with username in query string
	tampered := strings.Replace(initData, "john", "mallory", 1)
	_, err := VerifyInitData(tampered, botToken)
	if err == nil {
		t.Fatalf("expected error for tampered initData")
	}
}
