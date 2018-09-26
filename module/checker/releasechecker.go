package checker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/shurcooL/githubql"

	"github.com/github-releases-notifier/module/model"
)

// Checker has a githubql client to run queries and also knows about
// the current repositories releases to compare against.
type Checker struct {
	Logger   log.Logger
	Client   *githubql.Client
	releases map[string]model.Repository
}

// Run the queries and comparisons for the given repositories in a given interval.
func (c *Checker) Run(interval time.Duration, repositories []string, releases chan<- model.Repository, tag bool) {
	if c.releases == nil {
		c.releases = make(map[string]model.Repository)
	}

	for {
		for _, repoName := range repositories {
			s := strings.Split(repoName, "/")
			owner, name := s[0], s[1]

			nextRepo, err := c.query(owner, name)
			if true == tag {
				nextRepo, err = ApiV3tagChecker(owner, name)
				if nextRepo != (model.Repository{}) {
					releases <- nextRepo
					c.releases[repoName] = nextRepo
				}
			}

			if err != nil {
				level.Warn(c.Logger).Log(
					"msg", "failed to query the repository's releases",
					"owner", owner,
					"name", name,
					"err", err,
				)
				continue
			}

			// For debugging uncomment this next line
			//releases <- nextRepo

			currRepo, ok := c.releases[repoName]

			// We've queried the repository for the first time.
			// Saving the current state to compare with the next iteration.
			if !ok {
				c.releases[repoName] = nextRepo
				continue
			}

			if nextRepo.Release.PublishedAt.After(currRepo.Release.PublishedAt) {
				releases <- nextRepo
				c.releases[repoName] = nextRepo
			} else {
				level.Debug(c.Logger).Log(
					"msg", "no new release for repository",
					"owner", owner,
					"name", name,
				)
			}
		}
		time.Sleep(interval)
	}
}

// This should be improved in the future to make batch requests for all watched repositories at once
// TODO: https://github.com/shurcooL/githubql/issues/17

func (c *Checker) query(owner, name string) (model.Repository, error) {
	var query struct {
		Repository struct {
			ID          githubql.ID
			Name        githubql.String
			Description githubql.String
			URL         githubql.URI

			Releases struct {
				Edges []struct {
					Node struct {
						ID          githubql.ID
						Name        githubql.String
						Description githubql.String
						URL         githubql.URI
						PublishedAt githubql.DateTime
					}
				}
			} `graphql:"releases(last: 1)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubql.String(owner),
		"name":  githubql.String(name),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.Client.Query(ctx, &query, variables); err != nil {
		return model.Repository{}, err
	}

	repositoryID, ok := query.Repository.ID.(string)
	if !ok {
		return model.Repository{}, fmt.Errorf("can't convert repository id to string: %v", query.Repository.ID)
	}

	if len(query.Repository.Releases.Edges) == 0 {
		return model.Repository{}, fmt.Errorf("can't find any releases for %s/%s", owner, name)
	}
	latestRelease := query.Repository.Releases.Edges[0].Node

	releaseID, ok := latestRelease.ID.(string)
	if !ok {
		return model.Repository{}, fmt.Errorf("can't convert release id to string: %v", query.Repository.ID)
	}

	return model.Repository{
		ID:          repositoryID,
		Name:        string(query.Repository.Name),
		Owner:       owner,
		Description: string(query.Repository.Description),
		URL:         *query.Repository.URL.URL,

		Release: model.Release{
			ID:          releaseID,
			Name:        string(latestRelease.Name),
			Description: string(latestRelease.Description),
			URL:         *latestRelease.URL.URL,
			PublishedAt: latestRelease.PublishedAt.Time,
		},
	}, nil
}
