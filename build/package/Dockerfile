# Test stage
FROM golang:1.22 AS test

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY Makefile ./

RUN make test

# Build stage
FROM test AS build

# Accept GOOS and GOARCH as optional build arguments, which may be useful if you want to use this Dockerfile to build a
# binary for your host machine.
ARG GOOS=linux
ARG GOARCH=amd64

RUN make build GOOS=$GOOS GOARCH=$GOARCH

# Run stage
FROM scratch as run

COPY --from=build /app/bin/risk-api-linux-amd64 /risk-api

ENTRYPOINT ["/risk-api"]
