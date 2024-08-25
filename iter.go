package lfmap

import "github.com/benbjohnson/immutable"

func iterMap[K comparable, V any](m *immutable.Map[K, V], yield func(K, V) bool) {
	itr := m.Iterator()
	for !itr.Done() {
		k, v, _ := itr.Next()
		if !yield(k, v) {
			return
		}
	}
}
