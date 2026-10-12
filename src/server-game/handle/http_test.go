package handle

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

func TestCommandRequests(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		want       int
	}{
		{"delete rejects invalid input type", `{"action":"DeleteAccount","in":[]}`, http.StatusBadRequest},
		{"add bot rejects invalid input type", `{"action":"AddBot","in":[]}`, http.StatusBadRequest},
		{"delete bot rejects invalid input type", `{"action":"DeleteBot","in":[]}`, http.StatusBadRequest},
		{"empty email lists all", `{"action":"GetAccountList","in":{}}`, http.StatusOK},
		{"unvalidated email", `{"action":"GetAccountList","in":{"user_email":"invalid"}}`, http.StatusOK},
		{"valid list", `{"action":"GetAccountList","in":{"user_email":"nobody@example.invalid"}}`, http.StatusOK},
		{"create rejects invalid input type", `{"action":"CreateAccount","in":[]}`, http.StatusBadRequest},
		{"unknown command", `{"action":"Unknown"}`, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
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
		response := httptest.NewRecorder()
		HTTPHandle.ServeHTTP(response, request)
		if response.Code != http.StatusNotFound {
			t.Fatalf("legacy %s %s status = %d, want 404", route.method, route.path, response.Code)
		}
	}
}

func TestCreateAccountCommandResponse(t *testing.T) {
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
