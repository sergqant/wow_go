FROM golang:1.25-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0

COPY . .

RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
