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
	m.c.Update(func(oldValue *immutable.Map[K, V]) *immutable.Map[K, V] {
		t := tx[K, V]{oldValue}
		txn(&t)
		return t.m
	})
}

func (m *Map[K, V]) Clear() {
	m.Transaction(func(t Tx[K, V]) {
		t.Clear()
	})
}

func (m *Map[K, V]) Delete(key K) {
	m.Transaction(func(t Tx[K, V]) {
		t.Delete(key)
	})
}

func (m *Map[K, V]) Set(key K, value V) {
	m.Transaction(func(t Tx[K, V]) {
		t.Set(key, value)
	})
}

func (m *Map[K, V]) Range(yield func(key K, value V) bool) {
	iterMap(m.c.Value(), yield)
}
