# syntax=docker/dockerfile:1
# parent image for golang: https://hub.docker.com/_/golang
FROM golang:1.18.3

WORKDIR /app/api

WORKDIR /app/api/config
VOLUME ["/app/api/config"]

WORKDIR /app/api/databaseQueries
COPY databaseQueries/* ./
RUN ls -l

RUN go mod download && go mod verify && go mod tidy

WORKDIR /app/api/main
COPY main/* ./
RUN ls -l
# COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy
#RUN go build -v -o /usr/local/bin/app ./...
#CMD ["app"]

CMD ["go", "run", "main.go"]
