package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"rest-fiber/internal/enums"
)

func TestUserProfileAndAdminAccess(t *testing.T) {
	app, db, _, cleanup := setupTestApp(t)
	defer cleanup()

	migrateTestDB(t, db)

	adminEmail := randomEmail()
	adminPassword := "secret123"
	registerUser(t, app, "Admin User", adminEmail, adminPassword)
	adminUser := markUserVerifiedAndRole(t, db, adminEmail, enums.Admin)
	accessToken, _ := loginUser(t, app, adminEmail, adminPassword)

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	resp := requestJSON(t, app, http.MethodGet, "/api/v1/protected/users/me", nil, headers, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("users/me expected 200, got %d", resp.StatusCode)
	}
	body := decodeResponse(t, resp)
	if len(body.Data) == 0 {
		t.Fatal("users/me missing data")
	}

	userEmail := randomEmail()
	registerUser(t, app, "Member User", userEmail, "secret123")
	memberUser := markUserVerifiedAndRole(t, db, userEmail, enums.Member)
	if memberUser.ID == adminUser.ID {
		t.Fatal("member user should differ from admin user")
	}

	resp = requestJSON(t, app, http.MethodGet, "/api/v1/protected/users?page=1&limit=10", nil, headers, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("users list expected 200, got %d", resp.StatusCode)
	}
	body = decodeResponse(t, resp)
	var listData []any
	if err := json.Unmarshal(body.Data, &listData); err != nil {
		t.Fatalf("decode users list: %v", err)
	}
	if len(listData) == 0 {
		t.Fatal("users list missing data")
	}

	resp = requestJSON(t, app, http.MethodGet, "/api/v1/protected/users/"+memberUser.ID.String(), nil, headers, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("users/:id expected 200, got %d", resp.StatusCode)
	}
	body = decodeResponse(t, resp)
	if len(body.Data) == 0 {
		t.Fatal("users/:id missing data")
	}
}
