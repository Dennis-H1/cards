package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testConfig() Config {
	return Config{Username: "dennis", Password: "hunter2", SessionSecret: []byte("test-secret")}
}

func sessionCookie(t *testing.T, cfg Config) *http.Cookie {
	t.Helper()
	rec := httptest.NewRecorder()
	issueSession(rec, cfg)
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	return cookies[0]
}

func TestValidSessionRoundTrip(t *testing.T) {
	cfg := testConfig()
	cookie := sessionCookie(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	if !validSession(req, cfg) {
		t.Fatal("expected freshly issued session to be valid")
	}
}

func TestValidSessionRejectsTamperedCookie(t *testing.T) {
	cfg := testConfig()
	cookie := sessionCookie(t, cfg)
	cookie.Value = strings.Replace(cookie.Value, cookie.Value[:1], "x", 1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	if validSession(req, cfg) {
		t.Fatal("expected tampered cookie to be rejected")
	}
}

func TestValidSessionRejectsWrongSecret(t *testing.T) {
	cfg := testConfig()
	cookie := sessionCookie(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	otherCfg := cfg
	otherCfg.SessionSecret = []byte("different-secret")
	if validSession(req, otherCfg) {
		t.Fatal("expected cookie signed with a different secret to be rejected")
	}
}

func TestValidSessionRejectsExpiredCookie(t *testing.T) {
	cfg := testConfig()
	expiredPayload := "1" // unix epoch second 1, long past
	value := expiredPayload + "." + sign(cfg.SessionSecret, expiredPayload)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: value})
	if validSession(req, cfg) {
		t.Fatal("expected expired session to be rejected")
	}
}

func TestValidSessionRejectsMissingCookie(t *testing.T) {
	cfg := testConfig()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if validSession(req, cfg) {
		t.Fatal("expected request without a cookie to be rejected")
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	cfg := testConfig()
	body, _ := json.Marshal(loginRequest{Username: cfg.Username, Password: cfg.Password})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	LoginHandler(cfg)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if len(rec.Result().Cookies()) != 1 {
		t.Fatal("expected a session cookie to be set")
	}
}

func TestLoginHandlerWrongCredentials(t *testing.T) {
	cfg := testConfig()
	body, _ := json.Marshal(loginRequest{Username: cfg.Username, Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	LoginHandler(cfg)(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("expected no cookie to be set on failed login")
	}
}

func TestLogoutHandlerClearsCookie(t *testing.T) {
	cfg := testConfig()
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	rec := httptest.NewRecorder()

	LogoutHandler(cfg)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].MaxAge >= 0 {
		t.Fatalf("expected an expiring cookie, got %+v", cookies)
	}
}

func TestRequireAuthBlocksWithoutSession(t *testing.T) {
	cfg := testConfig()
	called := false
	handler := RequireAuth(cfg)(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if called {
		t.Fatal("expected wrapped handler not to be called")
	}
}

func TestRequireAuthAllowsWithSession(t *testing.T) {
	cfg := testConfig()
	cookie := sessionCookie(t, cfg)
	called := false
	handler := RequireAuth(cfg)(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !called {
		t.Fatal("expected wrapped handler to be called")
	}
}
