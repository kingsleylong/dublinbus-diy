# syntax=docker/dockerfile:1
# parent image for golang: https://hub.docker.com/_/golang
FROM "golang:1.18.3"

WORKDIR /usr/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .
#RUN go build -v -o /usr/local/bin/app ./...
#CMD ["app"]

CMD ["go", "version"]
