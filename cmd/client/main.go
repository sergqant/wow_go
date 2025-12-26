package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func solve(ctx context.Context, challenge string, difficulty int) (string, error) {
	target := strings.Repeat("0", difficulty)

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		nonce := strconv.Itoa(i)
		sum := sha256.Sum256([]byte(challenge + nonce))
		if strings.HasPrefix(hex.EncodeToString(sum[:]), target) {
			return nonce, nil
		}
	}
}

func readLine(ctx context.Context, r *bufio.Reader) (string, error) {
	type result struct {
		line string
		err  error
	}

	ch := make(chan result, 1)

	go func() {
		line, err := r.ReadString('\n')
		ch <- result{line, err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-ch:
		return res.line, res.err
	}
}

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}

	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		log.Fatal("connect:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// читаем challenge
	line, err := readLine(context.Background(), reader)
	if err != nil {
		log.Fatal("read challenge:", err)
	}

	var challenge string
	var diff int
	var timeoutMs int

	if _, err := fmt.Sscanf(
		line,
		"CHALLENGE %s %d %d",
		&challenge,
		&diff,
		&timeoutMs,
	); err != nil {
		log.Fatal("invalid challenge format:", err)
	}

	if diff <= 0 || timeoutMs <= 0 {
		log.Fatal("invalid challenge params")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutMs)*time.Millisecond,
	)
	defer cancel()

	// решаем pow
	solutionCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		nonce, err := solve(ctx, challenge, diff)
		if err != nil {
			errCh <- err
			return
		}

		solutionCh <- nonce
	}()

	// ждем ответ от сервера
	select {
	case nonce := <-solutionCh:
		// send solution
		if _, err := fmt.Fprintf(conn, "SOLUTION %s\n", nonce); err != nil {
			log.Fatal("send solution:", err)
		}

		resp, err := readLine(ctx, reader)
		if err != nil {
			log.Fatal("read quote:", err)
		}

		fmt.Print(resp)

	case err := <-errCh:
		log.Fatal("pow failed:", err)

	case <-ctx.Done():
		// server timeout or local timeout
		resp, err := readLine(context.Background(), reader)
		if err == nil {
			log.Fatal("server response:", strings.TrimSpace(resp))
		}
		log.Fatal(errors.New("pow timeout"))
	}
}
