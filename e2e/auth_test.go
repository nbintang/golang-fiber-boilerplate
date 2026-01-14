package e2e

import (
	"net/http"
	"testing"

	"rest-fiber/internal/enums"
)

func TestAuthRegisterLoginRefreshLogout(t *testing.T) {
	app, db, _, cleanup := setupTestApp(t)
	defer cleanup()

	migrateTestDB(t, db)

	email := randomEmail()
	password := "secret123"

	registerUser(t, app, "Auth User", email, password)
	markUserVerifiedAndRole(t, db, email, enums.Member)

	_, refreshCookie := loginUser(t, app, email, password)

	resp := requestJSON(t, app, http.MethodPost, "/api/v1/auth/refresh-token", nil, nil, []*http.Cookie{refreshCookie})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh-token expected 200, got %d", resp.StatusCode)
	}
	body := decodeResponse(t, resp)
	_ = extractToken(t, body.Data)
	newRefreshCookie := firstCookie(resp.Cookies(), "refresh_token")
	if newRefreshCookie == nil {
		t.Fatal("missing refresh_token cookie from refresh-token")
	}

	resp = requestJSON(t, app, http.MethodDelete, "/api/v1/auth/logout", nil, nil, []*http.Cookie{newRefreshCookie})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("logout expected 400, got %d", resp.StatusCode)
	}
}
