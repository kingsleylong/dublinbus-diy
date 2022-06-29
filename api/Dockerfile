# syntax=docker/dockerfile:1
# parent image for golang: https://hub.docker.com/_/golang
FROM golang:1.18.3

WORKDIR /app/api

WORKDIR /app/api/config
VOLUME ["/app/api/config"]

WORKDIR /app/api/databaseQueries
COPY ./databaseQueries/* ./app/api/databaseQueries
RUN ls -l
# pre-copy/cache go.mod for pre-downloading dependencies and only
# redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy

WORKDIR /app/api/main
COPY ./main/* /app/api/main
RUN ls -l
# COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy
#RUN go build -v -o /usr/local/bin/app ./...
#CMD ["app"]

CMD ["go", "run", "main.go"]
