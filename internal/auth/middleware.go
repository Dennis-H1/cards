package auth

import "net/http"

// RequireAuth returns a decorator that rejects requests without a valid
// session cookie before they reach next.
func RequireAuth(cfg Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !validSession(r, cfg) {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next(w, r)
		}
	}
}
