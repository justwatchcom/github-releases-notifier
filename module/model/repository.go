package model

import "net/url"

// Repository on GitHub.
type Repository struct {
	ID          string
	Name        string
	Owner       string
	Description string
	URL         url.URL
	Release     Release
}
