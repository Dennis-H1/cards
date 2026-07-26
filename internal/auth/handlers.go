package auth

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler checks the submitted credentials against cfg and, on success,
// issues a session cookie. It never reveals whether the username or the
// password was the wrong one.
func LoginHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request")
			return
		}

		validUsername := subtle.ConstantTimeCompare([]byte(req.Username), []byte(cfg.Username)) == 1
		validPassword := subtle.ConstantTimeCompare([]byte(req.Password), []byte(cfg.Password)) == 1
		if !validUsername || !validPassword {
			writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		issueSession(w, cfg)
		w.WriteHeader(http.StatusNoContent)
	}
}

// LogoutHandler clears the session cookie. It doesn't require an existing
// session -- clearing a cookie that isn't there is harmless.
func LogoutHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clearSession(w)
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
