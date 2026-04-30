package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Level     string    `json:"level"`
	Secret    string    `json:"secret"`
	Message   string    `json:"message"`
	ExpiresIn string    `json:"expires_in"`
	Timestamp time.Time `json:"timestamp"`
}

// WebhookNotifier sends alert payloads to an HTTP endpoint.
type WebhookNotifier struct {
	URL     string
	Client  *http.Client
	Headers map[string]string
}

// NewWebhookNotifier creates a WebhookNotifier with a default HTTP client.
func NewWebhookNotifier(url string, headers map[string]string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Headers: headers,
	}
}

// Send marshals the payload and POSTs it to the configured webhook URL.
func (w *WebhookNotifier) Send(payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.Headers {
		req.Header.Set(k, v)
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		if len(respBody) > 0 {
			return fmt.Errorf("webhook: unexpected status %d from %s: %s", resp.StatusCode, w.URL, respBody)
		}
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}
	return nil
}
