FROM golang:1.10-alpine3.7 as builder

ENV CGO_ENABLED=0

RUN apk --no-cache add make git
WORKDIR /go/src/github.com/justwatchcom/github-releases-notifier
ADD . /go/src/github.com/justwatchcom/github-releases-notifier
RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/justwatchcom/github-releases-notifier/github-releases-notifier-linux-amd64 /bin/github-releases-notifier
ENTRYPOINT [ "/bin/github-releases-notifier" ]
