package controller

import (
	"bytes"
	"fmt"
	"net/http"
)

type Notifier interface {
	Notify(message string) error
}

type DiscordNotifier struct {
	WebhookURL string
}

func (d *DiscordNotifier) Notify(message string) error {
	payload := fmt.Sprintf(`{"content": "%s"}`, message)
	return postMessage(d.WebhookURL, payload)
}

type TeamsNotifier struct {
	WebhookURL string
}

func (t *TeamsNotifier) Notify(message string) error {
	payload := fmt.Sprintf(`{
        "@type": "MessageCard",
        "@context": "http://schema.org/extensions",
        "summary": "Pod Restart Notification",
        "themeColor": "0078D7",
        "text": "%s"
    }`, message)
	return postMessage(t.WebhookURL, payload)
}

type SlackNotifier struct {
	WebhookURL string
}

func (s *SlackNotifier) Notify(message string) error {
	payload := fmt.Sprintf(`{"text": "%s"}`, message)
	return postMessage(s.WebhookURL, payload)
}

func postMessage(webhookURL string, payload string) error {
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("failed to send post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}

func sendNotification(message string, notifiers ...Notifier) error {
	var lastError error
	for _, notifier := range notifiers {
		if err := notifier.Notify(message); err != nil {
			lastError = err
		}
	}
	return lastError
}
