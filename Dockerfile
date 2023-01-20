# TODO: use multistage-build

FROM golang:1.19.5-alpine3.17
WORKDIR /app

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o ./main

CMD [ "./main" ]
