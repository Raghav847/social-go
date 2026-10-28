package mailer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"text/template"
)

type mailtrapClient struct {
	fromEmail  string
	apiKey     string
	httpClient *http.Client
	endpoint   string
}

func NewMailtrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("MAILTRAP_API_TOKEN is required")
	}
	return mailtrapClient{
		fromEmail:  fromEmail,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		endpoint:   "https://send.api.mailtrap.io/api/send",
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any) (int, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	type address struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}
	payload := struct {
		From    address   `json:"from"`
		To      []address `json:"to"`
		Subject string    `json:"subject"`
		HTML    string    `json:"html"`
	}{
		From:    address{Email: m.fromEmail},
		To:      []address{{Email: email, Name: username}},
		Subject: strings.TrimSpace(subject.String()),
		HTML:    body.String(),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequest(http.MethodPost, m.endpoint, bytes.NewReader(encoded))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
		if err != nil {
			return resp.StatusCode, fmt.Errorf("mailtrap API returned status %d: reading response: %w", resp.StatusCode, err)
		}
		return resp.StatusCode, fmt.Errorf("mailtrap API returned status %d: %s", resp.StatusCode, responseBody)
	}
	return resp.StatusCode, nil
}
