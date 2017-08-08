package main

import (
	"net/url"
	"time"
)

// Release of a repository tagged via GitHub.
type Release struct {
	ID          string
	Name        string
	Description string
	URL         url.URL
	PublishedAt time.Time
}
