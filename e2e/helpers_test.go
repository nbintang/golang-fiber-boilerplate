package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"rest-fiber/config"
	"rest-fiber/internal"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type httpResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

type paginationResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
	Meta       any             `json:"meta"`
}

func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB, config.Env, func()) {
	t.Helper()
	_ = os.Setenv("APP_ENV", "local")

	var bootstrap *internal.Bootstrap
	var db *gorm.DB
	var env config.Env

	app := fx.New(
		config.Module,
		infra.Module,
		internal.Module,
		fx.Populate(&bootstrap, &db, &env),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	if err := app.Start(ctx); err != nil {
		cancel()
		t.Fatalf("start app: %v", err)
	}
	cancel()

	cleanup := func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		_ = app.Stop(stopCtx)
	}

	return bootstrap.App, db, env, cleanup
}

func migrateTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		t.Fatalf("create pgcrypto extension: %v", err)
	}
	if err := db.Exec(`
	DO $$ BEGIN
		CREATE TYPE role_type AS ENUM ('ADMIN', 'MEMBER');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`).Error; err != nil {
		t.Fatalf("create role_type enum: %v", err)
	}
	if err := db.AutoMigrate(&user.User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
}

func registerUser(t *testing.T, app *fiber.App, name, email, password string) {
	t.Helper()
	payload := map[string]any{
		"name":       name,
		"email":      email,
		"avatar_url": "https://example.com/avatar.png",
		"password":   password,
	}
	resp := requestJSON(t, app, http.MethodPost, "/api/v1/auth/register", payload, nil, nil)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register expected 201, got %d", resp.StatusCode)
	}
}

func markUserVerifiedAndRole(t *testing.T, db *gorm.DB, email string, role enums.EUserRoleType) user.User {
	t.Helper()
	var u user.User
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		t.Fatalf("find user: %v", err)
	}
	if err := db.Model(&u).Updates(map[string]any{
		"is_email_verified": true,
		"role":              role,
	}).Error; err != nil {
		t.Fatalf("update user: %v", err)
	}
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	return u
}

func loginUser(t *testing.T, app *fiber.App, email, password string) (string, *http.Cookie) {
	t.Helper()
	payload := map[string]any{
		"email":    email,
		"password": password,
	}
	resp := requestJSON(t, app, http.MethodPost, "/api/v1/auth/login", payload, nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login expected 200, got %d", resp.StatusCode)
	}
	body := decodeResponse(t, resp)
	accessToken := extractToken(t, body.Data)
	refreshCookie := firstCookie(resp.Cookies(), "refresh_token")
	if refreshCookie == nil {
		t.Fatal("missing refresh_token cookie")
	}
	return accessToken, refreshCookie
}

func requestJSON(
	t *testing.T,
	app *fiber.App,
	method string,
	path string,
	body any,
	headers map[string]string,
	cookies []*http.Cookie,
) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode request: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func decodeResponse(t *testing.T, resp *http.Response) httpResponse {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	var out httpResponse
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return out
}

func extractToken(t *testing.T, data json.RawMessage) string {
	t.Helper()
	var payload map[string]string
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("decode token payload: %v", err)
	}
	token := payload["access_token"]
	if token == "" {
		t.Fatal("missing access_token")
	}
	return token
}

func firstCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func randomEmail() string {
	return "test+" + uuid.NewString() + "@example.com"
}
