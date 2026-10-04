package handle

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func TestCommandAuthorization(t *testing.T) {
	token := strings.Repeat("a", 32)
	previous := conf.ServerEnv.GameAPIToken
	conf.ServerEnv.GameAPIToken = token
	t.Cleanup(func() {
		conf.ServerEnv.GameAPIToken = previous
	})

	for _, tc := range []struct {
		name, token, body string
		want              int
	}{
		{"no token", "", `{"action":"GetAccountList"}`, http.StatusUnauthorized},
		{"wrong token", strings.Repeat("x", 32), `{"action":"GetAccountList"}`, http.StatusUnauthorized},
		{"delete rejects invalid input type", token, `{"action":"DeleteAccount","in":[]}`, http.StatusBadRequest},
		{"add bot rejects invalid input type", token, `{"action":"AddBot","in":[]}`, http.StatusBadRequest},
		{"delete bot rejects invalid input type", token, `{"action":"DeleteBot","in":[]}`, http.StatusBadRequest},
		{"empty email lists all", token, `{"action":"GetAccountList","in":{}}`, http.StatusOK},
		{"unvalidated email", token, `{"action":"GetAccountList","in":{"user_email":"invalid"}}`, http.StatusOK},
		{"valid list", token, `{"action":"GetAccountList","in":{"user_email":"nobody@example.invalid"}}`, http.StatusOK},
		{"create rejects invalid input type", token, `{"action":"CreateAccount","in":[]}`, http.StatusBadRequest},
		{"unknown command", token, `{"action":"Unknown"}`, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				request.Header.Set("Authorization", "Bearer "+tc.token)
			}
			response := httptest.NewRecorder()
			HTTPHandle.ServeHTTP(response, request)
			if response.Code != tc.want {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tc.want, response.Body.String())
			}
		})
	}

	for _, route := range []struct{ method, path string }{
		{http.MethodPost, "/api/accounts"},
		{http.MethodGet, "/api/accounts"},
		{http.MethodDelete, "/api/accounts/1"},
		{http.MethodPost, "/api/bots"},
		{http.MethodDelete, "/api/bots"},
	} {
		request := httptest.NewRequest(route.method, route.path, strings.NewReader(`{}`))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		HTTPHandle.ServeHTTP(response, request)
		if response.Code != http.StatusNotFound {
			t.Fatalf("legacy %s %s status = %d, want 404", route.method, route.path, response.Code)
		}
	}
}

func TestCreateAccountCommandResponse(t *testing.T) {
	token := strings.Repeat("a", 32)
	previous := conf.ServerEnv.GameAPIToken
	conf.ServerEnv.GameAPIToken = token
	t.Cleanup(func() { conf.ServerEnv.GameAPIToken = previous })

	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("t%x", suffix)
	email := name + "@example.invalid"
	t.Cleanup(func() {
		if err := model.DB.Where("name = ? AND user_email = ?", name, email).Delete(&model.Account{}).Error; err != nil {
			t.Errorf("delete test account: %v", err)
		}
	})

	body, err := json.Marshal(map[string]any{
		"action": "CreateAccount",
		"in":     map[string]string{"name": name, "password": "secret", "user_email": email},
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	HTTPHandle.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", response.Code, response.Body.String())
	}

	var result struct {
		Action string         `json:"action"`
		Out    map[string]any `json:"out"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Action != "CreateAccount" || len(result.Out) != 2 || result.Out["name"] != name || result.Out["user_email"] != email {
		t.Fatalf("unexpected command response: %s", response.Body.String())
	}
}

func TestValidAPIToken(t *testing.T) {
	token := strings.Repeat("a", 32)
	for _, tc := range []struct {
		header, configured string
		want               bool
	}{
		{"Bearer " + token, token, true},
		{"Bearer " + token, "", false},
		{"Bearer " + token, strings.Repeat("a", 31), false},
		{"Basic " + token, token, false},
		{"Bearer " + token + " extra", token, false},
	} {
		if got := validAPIToken(tc.header, tc.configured); got != tc.want {
			t.Errorf("validAPIToken(%q, configured length %d) = %v, want %v", tc.header, len(tc.configured), got, tc.want)
		}
	}
}
