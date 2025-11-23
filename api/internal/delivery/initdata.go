package delivery

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strings"
)

// VerifyInitData verifies Telegram Mini App initData according to Telegram docs.
// Returns parsed key->value map if valid.
func VerifyInitData(initData string, botToken string) (map[string]string, error) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}
	hash := vals.Get("hash")
	if hash == "" {
		return nil, errors.New("hash not found")
	}
	vals.Del("hash")
	keys := make([]string, 0, len(vals))
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := strings.Join(vals[k], "\n")
		parts = append(parts, k+"="+v)
	}
	dataCheckString := strings.Join(parts, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(hash)) {
		return nil, errors.New("invalid hash")
	}

	out := make(map[string]string)
	for k := range vals {
		out[k] = strings.Join(vals[k], "\n")
	}
	return out, nil
}
