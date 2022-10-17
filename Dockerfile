FROM golang:1.18 AS builder

WORKDIR /go/src/github.com/EquityZen/1password-datadog-events/

COPY . .

ARG gh_token
RUN echo "machine github.com login ezbuildbot password $gh_token" >> ~/.netrc
# Build app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o onepassword /go/src/github.com/EquityZen/1password-datadog-events/cmd

FROM alpine:latest
# vault AppRole
ARG secret_id
ENV SECRET_ID=$secret_id
ARG role_id
ENV ROLE_ID=$role_id

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# set working directory
WORKDIR /

# copy the binary from builder
COPY --from=builder /go/src/github.com/EquityZen/1password-datadog-events/cmd /usr/local/bin/onepassword

# run the binary
CMD ["onepassword"]