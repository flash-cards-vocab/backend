FROM golang:1.17.3-alpine as builder

LABEL maintainer="Aitugan Mirash"

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /api
WORKDIR /api

COPY . .
COPY .env .

RUN go mod tidy
RUN go get .
RUN go install .

# Build the Go api
RUN go build -o /build

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]
