package lfmap

type Tx[K comparable, V any] interface {
	Delete(key K)
	Get(key K) (value V, ok bool)
	Len() int
	Set(key K, value V)
	Range(yield func(key K, value V) bool)
}

