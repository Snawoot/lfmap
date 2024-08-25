package lfmap

import (
	"iter"

	"github.com/benbjohnson/immutable"
)

// Tx is a handle to map state usable from within transactions.
// It's methods are NOT safe for concurrent use.
type Tx[K comparable, V any] interface {
	// Clears the map.
	Clear()

	// Deletes the key.
	Delete(key K)

	// Returns the value for a given key and a flag indicating whether the key
	// exists. This flag distinguishes a nil value set on a key versus a
	// non-existent key in the map.
	Get(key K) (value V, ok bool)

	// Returns the number of elements in the map.
	Len() int

	// Map iterator suitable for use with range keyword.
	Range(yield func(key K, value V) bool)

	// Updates the map setting specified key to the new value.
	Set(key K, value V)
}

var _ Tx[string, string] = &tx[string, string]{}
var _ iter.Seq2[string, string] = (&tx[string, string]{}).Range

type tx[K comparable, V any] struct {
	m *immutable.Map[K, V]
}

func (t *tx[K, V]) Clear() {
	t.m = immutable.NewMap[K, V](newHasher[K]())
}

func (t *tx[K, V]) Delete(key K) {
	t.m = t.m.Delete(key)
}

func (t *tx[K, V]) Get(key K) (value V, ok bool) {
	return t.m.Get(key)
}

func (t *tx[K, V]) Len() int {
	return t.m.Len()
}

func (t *tx[K, V]) Set(key K, value V) {
	t.m = t.m.Set(key, value)
}

func (t *tx[K, V]) Range(yield func(key K, value V) bool) {
	iterMap(t.m, yield)
}
