package quotes

import (
	"crypto/rand"
	"math/big"
)

var data = []string{
	"Knowledge speaks, but wisdom listens.",
	"Wisdom begins in wonder.",
	"Turn your wounds into wisdom.",
	"The only true wisdom is knowing you know nothing.",
}

func Random() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data))))
	return data[n.Int64()]
}
