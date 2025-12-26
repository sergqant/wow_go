package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"task/word-of-wisdom/internal/pow"
	"task/word-of-wisdom/internal/server"
)

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return def
}

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}

	diff := envInt("POW_DIFFICULTY", 4)
	timeout := time.Duration(envInt("POW_TIMEOUT_MS", 5000)) * time.Millisecond

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	handler := &server.Handler{
		PoW:     &pow.SHA256PoW{Difficulty: diff},
		Timeout: timeout,
	}

	var wg sync.WaitGroup

	log.Println("Listening on", addr)

	go func() {
		<-ctx.Done()
		log.Println("Shutdown signal received")

		ln.Close() // прерываем Accept()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("Stopped accepting new connections")

				wg.Wait() // ждём активные соединения

				log.Println("Graceful shutdown complete")
				return
			default:
				continue
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			handler.Handle(ctx, conn)
		}()
	}
}
