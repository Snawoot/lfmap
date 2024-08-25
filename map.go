package lfmap

import (
	"github.com/Snawoot/occ"
	"github.com/benbjohnson/immutable"
)

type Map[K comparable, V any] struct {
	c occ.Container[immutable.Map[K, V]]
}

func New[K comparable, V any]() *Map[K, V] {
	m := new(Map[K, V])
	m.c.Update(func(_ *immutable.Map[K, V]) *immutable.Map[K, V] {
		return immutable.NewMap[K, V](newHasher[K]())
	})
	return m
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	return m.c.Value().Get(key)
}

func (m *Map[K, V]) Len(key K) int {
	return m.c.Value().Len()
}

func (m *Map[K, V]) Transaction(txn func(t Tx[K, V])) {
}
