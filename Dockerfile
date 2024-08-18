FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/main.go

FROM alpine:latest

ENV APP_ENV=production \
    APP_PORT=8080 \
    DB_PASSWORD=password

COPY --from=builder /go/bin/app /app
COPY configs /configs

RUN echo "APP_ENV=production" > .env && \
    echo "APP_PORT=8080" >> .env && \
    echo "DB_PASSWORD=qwerty" >> .env

EXPOSE 8080

CMD ["/app"]