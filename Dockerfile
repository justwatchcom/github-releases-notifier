FROM golang:1.10-alpine3.7 as builder

ADD . /go/src/github.com/justwatchcom/github-releases-notifier
WORKDIR /go/src/github.com/justwatchcom/github-releases-notifier

RUN apk --no-cache add make git

ENV CGO_ENABLED=0
RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/justwatchcom/github-releases-notifier /bin/
ENTRYPOINT [ "/bin/github-releases-notifier" ]
