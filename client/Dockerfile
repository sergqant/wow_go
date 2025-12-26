FROM golang:1.25-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0

COPY . .

RUN go build -o client ./cmd/client

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/client .
CMD ["./client"]
