package lfmap

import (
	"iter"

	"github.com/benbjohnson/immutable"
)

type Tx[K comparable, V any] interface {
	Delete(key K)
	Get(key K) (value V, ok bool)
	Len() int
	Range(yield func(key K, value V) bool)
	Set(key K, value V)
}

var _ Tx[string, string] = &tx[string, string]{}
var _ iter.Seq2[string, string] = Tx[string, string](nil).Range

type tx[K comparable, V any] struct {
	m *immutable.Map[K, V]
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
	itr := t.m.Iterator()
	for !itr.Done() {
		k, v, _ := itr.Next()
		if !yield(k, v) {
			return
		}
	}
}
