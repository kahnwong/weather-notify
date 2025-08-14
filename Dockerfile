FROM golang:1.25-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -ldflags "-w -s" -o /weather-notify

# hadolint ignore=DL3007
FROM alpine:latest AS build-release-stage

WORKDIR /

RUN apk add --no-cache tzdata

COPY --from=build-stage /weather-notify /weather-notify

RUN chmod +x /weather-notify

ENTRYPOINT ["/weather-notify"]
