package mailer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMailtrapSend(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusUnauthorized, http.StatusUnprocessableEntity, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Content-Type") != "application/json" {
					t.Error("incorrect request method or headers")
				}
				var payload struct {
					From          struct{ Email string }
					To            []struct{ Email, Name string }
					Subject, HTML string
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Error(err)
				}
				if payload.From.Email != "hello@demomailtrap.co" || len(payload.To) != 1 {
					t.Error("incorrect sender or recipients")
				} else if payload.To[0].Email != "user@example.com" || payload.To[0].Name != "test-user" {
					t.Error("incorrect recipient")
				}
				if payload.Subject != "Finish Registration with GopherSocial" || !strings.Contains(payload.HTML, "test-user") || !strings.Contains(payload.HTML, "https://example.com/confirm/token") {
					t.Error("incorrect rendered email")
				}
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"success":true}`))
			}))
			defer server.Close()
			client, err := NewMailtrapClient("test-token", "hello@demomailtrap.co")
			if err != nil {
				t.Fatal(err)
			}
			client.endpoint = server.URL
			client.httpClient = server.Client()
			data := struct{ UserName, ActivationURL string }{"test-user", "https://example.com/confirm/token"}
			gotStatus, err := client.Send(UserWelcomeTemplate, "test-user", "user@example.com", data)
			if gotStatus != status || (err != nil) != (status >= 300) {
				t.Fatalf("status = %d, error = %v", gotStatus, err)
			}
		})
	}
}

func TestMailtrapRequiresToken(t *testing.T) {
	if _, err := NewMailtrapClient("", "hello@demomailtrap.co"); err == nil {
		t.Fatal("expected missing token error")
	}
}
