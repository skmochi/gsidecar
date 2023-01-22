# build stage
FROM golang:1.19.5-alpine3.17 AS builder
WORKDIR /app

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -ldflags="-s -w" -trimpath -o ./main


# run stage
FROM alpine:3.17.1
WORKDIR /app

COPY --from=builder /app/main .
CMD [ "./main" ]
