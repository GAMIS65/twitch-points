package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func SendWebHook(text string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return errors.New("WEBHOOK_URL environment variable is not set")
	}

	payload := map[string]string{
		"content": text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errors.New("error encoding JSON: " + err.Error())
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("error creating request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("error sending webhook request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("failed to send webhook, status: " + resp.Status)
	}

	return nil
}
