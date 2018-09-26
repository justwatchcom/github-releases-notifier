package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SlackSender has the hook to send slack notifications.
type SlackSender struct {
	Hook string
}

type slackPayload struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
}

// Send a notification with a formatted message build from the repository.
func (s *SlackSender) Send(repository Repository) error {
	payload := slackPayload{
		Username:  "GitHub Releases",
		IconEmoji: ":github:",
		Text: fmt.Sprintf(
			"<%s|%s/%s>: <%s|%s> released",
			repository.URL.String(),
			repository.Owner,
			repository.Name,
			repository.Release.URL.String(),
			repository.Release.Name,
		),
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.Hook, bytes.NewReader(payloadData))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	req = req.WithContext(ctx)
	defer cancel()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("request didn't respond with 200 OK: %s, %s", resp.Status, body)
	}

	return nil
}
