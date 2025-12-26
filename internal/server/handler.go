package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net"
	"strings"
	"time"

	"task/word-of-wisdom/internal/pow"
	"task/word-of-wisdom/internal/protocol"
	"task/word-of-wisdom/internal/quotes"
)

type Handler struct {
	PoW     *pow.SHA256PoW
	Timeout time.Duration
}

func randomChallenge() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	return n.Text(16)
}

func (h *Handler) Handle(parentCtx context.Context, conn net.Conn) {
	defer conn.Close()

	ctx, cancel := context.WithTimeout(parentCtx, h.Timeout)
	defer cancel()

	if err := conn.SetDeadline(time.Now().Add(h.Timeout)); err != nil {
		log.Printf("set deadline error: %v", err)
		
		return
	}

	challenge := randomChallenge()
	if _, err := conn.Write([]byte(protocol.Challenge(
		challenge,
		h.PoW.Difficulty,
		int(h.Timeout.Milliseconds()),
	))); err != nil {
		log.Printf("write challenge error: %v", err)

		return
	}

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Printf("read error: %v", err)
		}

		return
	}

	select {
	case <-ctx.Done():
		_, err = conn.Write([]byte(protocol.Error("timeout")))
		if err != nil {
			log.Printf("write timeout error: %v", err)
		}
		return
	default:
	}

	parts := strings.Fields(line)
	if len(parts) != 2 || parts[0] != "SOLUTION" {
		if _, err := conn.Write([]byte(protocol.Error("bad format"))); err != nil {
			log.Printf("write error: %v", err)
		}

		return
	}

	if !h.PoW.Verify(challenge, parts[1]) {
		if _, err := conn.Write([]byte(protocol.Error("invalid pow"))); err != nil {
			log.Printf("write error: %v", err)
		}

		return
	}

	if _, err := conn.Write([]byte(protocol.OK(quotes.Random()))); err != nil {
		log.Printf("write quote error: %v", err)

		return
	}
}
