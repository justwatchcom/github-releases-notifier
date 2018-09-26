package checker

import (
	"context"
	"os"
	"fmt"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"net/url"
	"time"

	"github.com/github-releases-notifier/module/model"
	"github.com/github-releases-notifier/module/redis"
)

// Function for working with github api v3 and check if new tags are published
func ApiV3tagChecker(owner, name string) (model.Repository, error) {
	redis.СonnectToRedis()
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")

	if !ok {
		fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.ListOptions{}

	repo, _, err := client.Repositories.Get(ctx, owner, name)
	tags, _, err := client.Repositories.ListTags(ctx, owner, name, opt)

	if err != nil {
		return model.Repository{}, fmt.Errorf("cant find repo")
	}

	var LastTag github.RepositoryTag
	for _, tag := range tags {
		LastTag = *tag
		break
	}

	repourl, err := url.ParseRequestURI(repo.GetURL())
	if err != nil {
		return model.Repository{}, fmt.Errorf("cant convert url")
	}

	if redis.GetValue(name) == LastTag.GetName() {
		return model.Repository{}, nil
	}

	redis.SetKey(name, LastTag.GetName())

	return model.Repository{
		ID:          string(repo.GetID()),
		Name:        string(repo.GetName()),
		Owner:       owner,
		Description: string(repo.GetDescription()),
		URL:         *repourl,

		Release: model.Release{
			ID:			 "1",
			Name:        LastTag.GetName(),
			Description: "It is Tag",
			URL:         *repourl,
			PublishedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}, nil
}