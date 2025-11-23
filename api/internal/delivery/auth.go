package delivery

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"timerwallet-api/internal/db"

	"github.com/golang-jwt/jwt/v5"
)

// TelegramAuthHandler validates initData, creates/loads user, issues tokens and sets cookies
func TelegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		InitData string `json:"initData"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	vals, err := url.ParseQuery(body.InitData)
	if err != nil {
		http.Error(w, "invalid initData format", http.StatusBadRequest)
		return
	}
	// verify
	hash := vals.Get("hash")
	if hash == "" {
		http.Error(w, "hash not found", http.StatusBadRequest)
		return
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

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}
	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(hash)) {
		http.Error(w, "invalid initData hash", http.StatusUnauthorized)
		return
	}

	// extract fields
	tidStr := vals.Get("id")
	if tidStr == "" {
		http.Error(w, "id not found", http.StatusBadRequest)
		return
	}
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	username := vals.Get("username")
	firstName := vals.Get("first_name")
	lastName := vals.Get("last_name")
	authDate := vals.Get("auth_date")

	// find or create
	var userID int64
	var createdAt sql.NullString
	var dbUsername, dbFirstName, dbLastName sql.NullString
	var dbAuthDate sql.NullString

	err = db.DB.QueryRow("SELECT id, username, first_name, last_name, auth_date, created_at FROM telegram_users WHERE telegram_id=$1", tid).Scan(&userID, &dbUsername, &dbFirstName, &dbLastName, &dbAuthDate, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.DB.QueryRow("INSERT INTO telegram_users (telegram_id, username, first_name, last_name, auth_date) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at", tid, username, firstName, lastName, authDate).Scan(&userID, &createdAt)
			if err != nil {
				log.Printf("failed to insert user: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			dbUsername = sql.NullString{String: username, Valid: username != ""}
			dbFirstName = sql.NullString{String: firstName, Valid: firstName != ""}
			dbLastName = sql.NullString{String: lastName, Valid: lastName != ""}
			dbAuthDate = sql.NullString{String: authDate, Valid: authDate != ""}
		} else {
			log.Printf("db error: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	accessClaims := jwt.MapClaims{"sub": strconv.FormatInt(userID, 10), "exp": accessExp.Unix(), "type": "access"}
	refreshClaims := jwt.MapClaims{"sub": strconv.FormatInt(userID, 10), "exp": refreshExp.Unix(), "type": "refresh"}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("token gen err: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("token gen err: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	cookieSecure := os.Getenv("COOKIE_SECURE") == "1"
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: accessToken, Path: "/", HttpOnly: true, Secure: cookieSecure, SameSite: http.SameSiteLaxMode, Expires: accessExp})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: refreshToken, Path: "/", HttpOnly: true, Secure: cookieSecure, SameSite: http.SameSiteLaxMode, Expires: refreshExp})

	userResp := map[string]interface{}{"id": userID, "first_name": dbFirstName.String, "last_name": dbLastName.String, "username": dbUsername.String, "roles": []string{"user"}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResp)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := cookie.Value
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	uid, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var id int64
	var username, firstName, lastName sql.NullString
	err = db.DB.QueryRow("SELECT id, username, first_name, last_name FROM telegram_users WHERE id=$1", uid).Scan(&id, &username, &firstName, &lastName)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	userResp := map[string]interface{}{"id": id, "first_name": firstName.String, "last_name": lastName.String, "username": username.String, "roles": []string{"user"}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResp)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := cookie.Value
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if tType, _ := claims["type"].(string); tType != "refresh" {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	uid, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	accessClaims := jwt.MapClaims{"sub": strconv.FormatInt(uid, 10), "exp": accessExp.Unix(), "type": "access"}
	refreshClaims := jwt.MapClaims{"sub": strconv.FormatInt(uid, 10), "exp": refreshExp.Unix(), "type": "refresh"}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	cookieSecure := os.Getenv("COOKIE_SECURE") == "1"
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: accessToken, Path: "/", HttpOnly: true, Secure: cookieSecure, SameSite: http.SameSiteLaxMode, Expires: accessExp})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: refreshToken, Path: "/", HttpOnly: true, Secure: cookieSecure, SameSite: http.SameSiteLaxMode, Expires: refreshExp})
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Path: "/", HttpOnly: true, Secure: false, MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", HttpOnly: true, Secure: false, MaxAge: -1})
	w.WriteHeader(http.StatusOK)
}
