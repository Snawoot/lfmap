package lfmap

import (
	"hash/maphash"
)

type hasher[K comparable] struct{
	seed maphash.Seed
}

func newHasher[K comparable]() hasher[K] {
	return hasher[K]{
		seed: maphash.MakeSeed(),
	}
}

func (h hasher[K]) Hash(key K) uint32 {
	return uint32(maphash.Comparable[K](h.seed, key))
}

func (h hasher[K]) Equal(a, b K) bool {
	return a == b
}
