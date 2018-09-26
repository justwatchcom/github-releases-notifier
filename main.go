package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubql"
	"golang.org/x/oauth2"

	"github.com/github-releases-notifier/module/model"
	"github.com/github-releases-notifier/module/sender"
	releasechecker "github.com/github-releases-notifier/module/checker"
)

// Config of env and args
type Config struct {
	GithubToken  string        `arg:"env:GITHUB_TOKEN"`
	Interval     time.Duration `arg:"env:INTERVAL"`
	LogLevel     string        `arg:"env:LOG_LEVEL"`
	Repositories []string      `arg:"-r,separate"`
	SlackHook    string        `arg:"env:SLACK_HOOK"`
	IsTagChecker bool 		   `arg:"env:TAG_CHECKER"`
}

// Token returns an oauth2 token or an error.
func (c Config) Token() *oauth2.Token {
	return &oauth2.Token{AccessToken: c.GithubToken}
}

func main() {
	_ = godotenv.Load()

	c := Config{
		Interval: time.Second,
		LogLevel: "info",
	}
	arg.MustParse(&c)

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.Caller(5),
	)

	level.SetKey("severity")
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	tokenSource := oauth2.StaticTokenSource(c.Token())
	client := oauth2.NewClient(context.Background(), tokenSource)
	checker := &releasechecker.Checker{
		Logger: logger,
		Client: githubql.NewClient(client),
	}

	releases := make(chan model.Repository)
	go checker.Run(c.Interval, c.Repositories, releases, c.IsTagChecker)

	slack := sender.SlackSender{Hook: c.SlackHook}

	level.Info(logger).Log("msg", "waiting for new releases")
	for repository := range releases {
		if err := slack.Send(repository); err != nil {
			level.Warn(logger).Log(
				"msg", "failed to send release to messenger",
				"err", err,
			)
			continue
		}
	}
}
