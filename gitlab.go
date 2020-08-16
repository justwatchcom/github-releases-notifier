package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// GitlabSender has the URL, ProjectID and Token to create GitLab issues.
type GitlabSender struct {
	Hostname  string
	APIToken  string
	ProjectID int
	Labels    string
	logger    log.Logger
}

type gitlabPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Labels      string `json:"labels"`
}

// Send a notification with a formatted message build from the repository.
func (g *GitlabSender) Send(repository Repository) error {
	payload := gitlabPayload{
		Title: fmt.Sprintf(
			":arrow_up: New version of %s released: %s",
			repository.Name,
			repository.Release.Name,
		),
		Description: fmt.Sprintf(
			"More info: %s",
			repository.Release.URL.String(),
		),
		Labels: g.Labels,
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	url := strings.Join([]string{"https:/", g.Hostname, "api/v4/projects", strconv.Itoa(g.ProjectID), "issues"}, "/")
	level.Debug(g.logger).Log(
		"msg", "attempting to post issue",
		"url", url,
	)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadData))
	if err != nil {
		return err
	}
	req.Header.Set("PRIVATE-TOKEN", g.APIToken)
	req.Header.Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	req = req.WithContext(ctx)
	defer cancel()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("request didn't respond with 201 Created: %s, %s", resp.Status, body)
	}

	return nil
}
