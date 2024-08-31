package lfmap

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"testing"
)

func TestSmoke(t *testing.T) {
	m := New[string, int]()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("k%d", i)
			value := i * 100
			m.Set(key, value)
		}(i)
	}
	wg.Wait()

	for k, v := range m.Range {
		t.Logf("key = %s, value = %d", k, v)
		if k != fmt.Sprintf("k%d", v/100) {
			t.Fail()
		}
	}
}

const MAP_LOAD = 100_000

func benchPar[T any](b *testing.B, factory func() T, f func(b *testing.B, m T)) {
	m := factory()
	b.ResetTimer()
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(b, m)
		}()
	}
	wg.Wait()
}

func BenchmarkLFMapSet(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] { return New[int, int]() },
		func(b *testing.B, m *Map[int, int]) {
			r := rand.New(&rand.PCG{})
			for i := 0; i < b.N; i++ {
				m.Set(r.Int(), i)
			}
		},
	)
}

func BenchmarkSyncMapSet(b *testing.B) {
	benchPar(b,
		func() *sync.Map { return &sync.Map{} },
		func(b *testing.B, m *sync.Map) {
			r := rand.New(&rand.PCG{})
			for i := 0; i < b.N; i++ {
				m.Store(r.Int(), i)
			}
		},
	)
}

func BenchmarkLFMapGet(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] {
			m := New[int, int]()
			for i := 0; i < MAP_LOAD; i++ {
				m.Set(i, i)
			}
			return m
		},
		func(b *testing.B, m *Map[int, int]) {
			r := rand.New(&rand.PCG{})
			for i := 0; i < b.N; i++ {
				_, _ = m.Get(r.IntN(MAP_LOAD))
			}
		},
	)
}

func BenchmarkSyncMapGet(b *testing.B) {
	benchPar(b,
		func() *sync.Map {
			m := &sync.Map{}
			for i := 0; i < MAP_LOAD; i++ {
				m.Store(i, i)
			}
			return m
		},
		func(b *testing.B, m *sync.Map) {
			r := rand.New(&rand.PCG{})
			for i := 0; i < b.N; i++ {
				_, _ = m.Load(r.IntN(MAP_LOAD))
			}
		},
	)
}

func BenchmarkLFMapRange100000(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] {
			m := New[int, int]()
			for i := 0; i < MAP_LOAD; i++ {
				m.Set(i, i)
			}
			return m
		},
		func(b *testing.B, m *Map[int, int]) {
			for i := 0; i < b.N; i++ {
				var ctr int
				for _, _ = range m.Range {
					ctr++
				}
				if ctr != MAP_LOAD {
					b.Fatalf("unexpected number of interations: %d", ctr)
				}
			}
		},
	)
}

func BenchmarkSyncMapRange100000(b *testing.B) {
	benchPar(b,
		func() *sync.Map {
			m := &sync.Map{}
			for i := 0; i < MAP_LOAD; i++ {
				m.Store(i, i)
			}
			return m
		},
		func(b *testing.B, m *sync.Map) {
			for i := 0; i < b.N; i++ {
				var ctr int
				for _, _ = range m.Range {
					ctr++
				}
				if ctr != MAP_LOAD {
					b.Fatalf("unexpected number of interations: %d", ctr)
				}
			}
		},
	)
}
