// Package auth gates the REST API behind a single-account login, separate
// from the MCP server's bearer-token auth. There is no users table -- the
// one account's credentials come from env vars, same pattern as MCP_API_KEY.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cookieName = "karteikarten_session"
	sessionTTL = 30 * 24 * time.Hour
)

// Config holds the single account's credentials and the key used to sign
// session cookies.
type Config struct {
	Username      string
	Password      string
	SessionSecret []byte
}

func issueSession(w http.ResponseWriter, cfg Config) {
	expiry := time.Now().Add(sessionTTL)
	payload := strconv.FormatInt(expiry.Unix(), 10)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    payload + "." + sign(cfg.SessionSecret, payload),
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func validSession(r *http.Request, cfg Config) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	payload, sig, ok := strings.Cut(c.Value, ".")
	if !ok {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(sig), []byte(sign(cfg.SessionSecret, payload))) != 1 {
		return false
	}
	expiry, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return false
	}
	return time.Now().Unix() < expiry
}

func sign(secret []byte, payload string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
