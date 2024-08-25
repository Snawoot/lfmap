package lfmap

import "testing"

func TestHasher(t *testing.T) {
	h := newHasher[string]()
	h1 := h.Hash("hello")
	h2 := h.Hash("hello")
	h3 := h.Hash("world")
	if h1 != h2 {
		t.Fatal("h1 != h2")
	}
	if h1 == h3 {
		t.Fatal("h2 == h3")
	}
}
