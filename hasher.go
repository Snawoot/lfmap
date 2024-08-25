package lfmap

import (
	"github.com/dolthub/maphash"
)

type hasher[K comparable] maphash.Hasher[K]

func newHasher[K comparable]() hasher[K] {
	return hasher[K](maphash.NewHasher[K]())
}

func (h hasher[K]) Hash(key K) uint32 {
	return uint32(maphash.Hasher[K](h).Hash(key))
}

func (h hasher[K]) Equal(a, b K) bool {
	return a == b
}
