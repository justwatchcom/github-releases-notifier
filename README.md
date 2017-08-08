# github-releases-notifier

[![Build Status](https://travis-ci.org/justwatchcom/github-releases-notifier.svg?branch=master)](https://travis-ci.org/justwatchcom/github-releases-notifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/justwatchcom/github-releases-notifier)](https://goreportcard.com/report/github.com/justwatchcom/github-releases-notifier)
[![Docker Pulls](https://img.shields.io/docker/pulls/justwatch/github-releases-notifier.svg?maxAge=604800)](https://hub.docker.com/r/justwatch/github-releases-notifier)

Receive Slack notifications if a new release of your favorite software is available on GitHub.

![screenshot.png](screenshot.png)

### Watching repositories

To watch repositories simply add them to the list of arguments `-r=kubernetes/kubernetes -r=prometheus/prometheus` and so on.

### Deploying

1. Get a URL to send WebHooks to your Slack from https://api.slack.com/incoming-webhooks.
2. Get a token for scraping GitHub: [https://help.github.com/](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line).

#### Docker

```
docker run --rm -e GITHUB_TOKEN=XXX -e SLACK_HOOK=https://hooks.slack.com/... justwatch/github-releases-notifier -r=kubernetes/kubernetes
```

#### docker-compose

1. Change into the `deployments/` folder.
2. Open `docker-compose.yml`
3. Change the token in the environment section to the ones obtained above.
4. `docker-compose up`

#### Kubernetes

```bash
kubectl create secret generic github-releases-notifier \
        --from-literal=github=XXX` \
        --from-literal=slack=XXX
```

After creating the secret with your credentials you can apply the deployment:

`kubectl apply -f deployments/kubernetes.yml`

That's it.