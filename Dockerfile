FROM golang:alpine as builder

RUN mkdir /goack
WORKDIR /goack
COPY go.mod .
COPY go.sum .
COPY .env .

RUN apk add --update --no-cache ca-certificates git

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
FROM alpine:latest
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache redis
RUN apk add --no-cache curl


RUN mkdir /app
WORKDIR /app

COPY --from=builder /goack/goack .
COPY --from=builder /goack/.env .
COPY --from=builder /goack/script.sh .

RUN ["chmod", "+x", "./script.sh"]
CMD "./script.sh"
