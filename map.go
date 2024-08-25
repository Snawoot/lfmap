// Package lfmap provides generic concurrent lock-free map implementation,
// using immutable map and atomic swaps.
package lfmap

import (
	"github.com/Snawoot/occ"
	"github.com/benbjohnson/immutable"
)

// Map is an instance of concurrent map. All Map methods are safe for concurrent use.
type Map[K comparable, V any] struct {
	c occ.Container[immutable.Map[K, V]]
}

// New returns a new instance of empty Map.
func New[K comparable, V any]() *Map[K, V] {
	m := new(Map[K, V])
	m.c.Update(func(_ *immutable.Map[K, V]) *immutable.Map[K, V] {
		return immutable.NewMap[K, V](newHasher[K]())
	})
	return m
}

// Get returns the value for a given key and a flag indicating whether the key
// exists. This flag distinguishes a nil value set on a key versus a
// non-existent key in the map.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	return m.c.Value().Get(key)
}

// Len returns the number of elements in the map.
func (m *Map[K, V]) Len(key K) int {
	return m.c.Value().Len()
}

// Transaction executes operations made by txn function, allowing complex
// interactions with map to be performed atomically and consistently.
//
// txn function can be invoked more than once in case if collision happened
// during update.
func (m *Map[K, V]) Transaction(txn func(t Tx[K, V])) {
	m.c.Update(func(oldValue *immutable.Map[K, V]) *immutable.Map[K, V] {
		t := tx[K, V]{oldValue}
		txn(&t)
		return t.m
	})
}

// Clear empties map.
// It's a shortcut for a transaction with just clear operation in it.
func (m *Map[K, V]) Clear() {
	m.Transaction(func(t Tx[K, V]) {
		t.Clear()
	})
}

// Delete updates the map removing specified key.
// It's a shortcut for a transaction with just delete operation in it.
func (m *Map[K, V]) Delete(key K) {
	m.Transaction(func(t Tx[K, V]) {
		t.Delete(key)
	})
}

// Set updates the map setting specified key to the new value. It's a shorcut
// for a [Map.Transaction] invoked with function setting only one value.
func (m *Map[K, V]) Set(key K, value V) {
	m.Transaction(func(t Tx[K, V]) {
		t.Set(key, value)
	})
}

// Map iterator suitable for use with range keyword.
func (m *Map[K, V]) Range(yield func(key K, value V) bool) {
	iterMap(m.c.Value(), yield)
}
