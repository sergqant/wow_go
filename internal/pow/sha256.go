package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type SHA256PoW struct {
	Difficulty int
}

func (p *SHA256PoW) Verify(challenge, nonce string) bool {
	sum := sha256.Sum256([]byte(challenge + nonce))
	hash := hex.EncodeToString(sum[:])
	return strings.HasPrefix(hash, strings.Repeat("0", p.Difficulty))
}
