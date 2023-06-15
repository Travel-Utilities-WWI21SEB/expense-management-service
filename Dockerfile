# Golang Base Image
FROM golang:1.20.5-alpine3.18 AS build

## Build the executable in the first stage

WORKDIR /go/app

COPY go.mod go.sum ./
COPY src/ ./src

RUN go mod download
RUN go build -o expense-api ./src

## Serve only the compiled binary in the second stage
FROM alpine:3.18.0 AS serve

## Neccessary to run a health check in our docker-compose file
RUN apk --update --no-cache add curl

# Copy the Pre-built binary file from the previous stage
COPY --from=build /go/app/expense-api /go/app/expense-api

CMD ["/go/app/expense-api"]
