# Test task for Server Engineer Golang

## Task description
Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.

## How to build
Server:

```
docker build -f server/Dockerfile -t wow-server .
```

Client:

```
docker build -f client/Dockerfile -t wow-client .
```


## How to run
Server:

```
docker run -p 8083:8080 \
  -e POW_DIFFICULTY=4 \
  -e POW_TIMEOUT_MS=5000 \
  -e PORT=":8080" \
  wow-server
```

Client:
```
docker run --rm \
  -e PORT="host.docker.internal:8083" \
  wow-client
```
