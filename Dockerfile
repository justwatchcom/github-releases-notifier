FROM golang:1.10-alpine3.7 as builder

ENV CGO_ENABLED=0
ARG DEP_VERSION=v0.4.1

RUN apk --no-cache add make git
ADD https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 /usr/local/bin
RUN mv /usr/local/bin/dep-linux-amd64 /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/justwatchcom/github-releases-notifier
ADD Makefile Gopkg.lock Gopkg.toml /go/src/github.com/justwatchcom/github-releases-notifier/
RUN make dep

ADD . /go/src/github.com/justwatchcom/github-releases-notifier
RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/justwatchcom/github-releases-notifier/github-releases-notifier-linux-amd64 /bin/github-releases-notifier
ENTRYPOINT [ "/bin/github-releases-notifier" ]
