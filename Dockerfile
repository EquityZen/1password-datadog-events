#
# Stage 1 - Compile go code
#

FROM golang:alpine as builder

# Install XZ and UPZ
RUN apk update apk add ca-certificates xz upx git && rm -rf /var/cache/apk/* \
                   && rm -rf /var/lib/apt/lists/*

ARG gh_token
RUN echo "machine github.com login ezbuildbot password $gh_token" >> ~/.netrc
ARG git_hash=""

# Create work dir
WORKDIR /go/src/github.com/EquityZen/1password-datadog-events/

# copy entire directory
COPY . .

# Build app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.GitCommit=$git_hash" -a -installsuffix cgo -o onepassword github.com/EquityZen/1password-datadog-events/cmd/

# strip and compress the binary
RUN upx onepassword

#
# Stage 2 -  Add complied binary to new clean container
#
FROM alpine

# vault AppRole
ARG secret_id
ENV SECRET_ID=$secret_id
ARG role_id
ENV ROLE_ID=$role_id

# add ca-certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# set working directory
WORKDIR /

# copy the binary from builder
COPY --from=builder /go/src/github.com/EquityZen/1password-datadog-events/cmd /usr/local/bin/onepassword

# run the binary
CMD ["onepassword"]