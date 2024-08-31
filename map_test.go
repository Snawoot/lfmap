package lfmap

import (
	"fmt"
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

const MAP_LOAD = 1000_000

// dummy output variable to ensure invoked function are not optimized out
var BenchDrain = make([]any, runtime.NumCPU())

func benchPar[T any](b *testing.B, factory func() T, f func(b *testing.B, m T, drain *any)) {
	m := factory()
	b.ResetTimer()
	var wg sync.WaitGroup
	for drainIdx := range runtime.NumCPU() {
		wg.Add(1)
		go func(drainIdx int) {
			defer wg.Done()
			f(b, m, &BenchDrain[drainIdx])
		}(drainIdx)
	}
	wg.Wait()
}

func BenchmarkLFMapSet(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] { return New[int, int]() },
		func(b *testing.B, m *Map[int, int], _ *any) {
			for i := 0; i < b.N; i++ {
				m.Set(i%MAP_LOAD, i)
			}
		},
	)
}

func BenchmarkSyncMapSet(b *testing.B) {
	benchPar(b,
		func() *sync.Map { return &sync.Map{} },
		func(b *testing.B, m *sync.Map, _ *any) {
			for i := 0; i < b.N; i++ {
				m.Store(i%MAP_LOAD, i)
			}
		},
	)
}

func BenchmarkLFMapGet(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] {
			m := New[int, int]()
			// test 50% key miss
			for i := 0; i < 2*MAP_LOAD; i += 2 {
				m.Set(i, i)
			}
			return m
		},
		func(b *testing.B, m *Map[int, int], drain *any) {
			for i := 0; i < b.N; i++ {
				*drain, _ = m.Get(i % (2 * MAP_LOAD))
			}
		},
	)
}

func BenchmarkSyncMapGet(b *testing.B) {
	benchPar(b,
		func() *sync.Map {
			m := &sync.Map{}
			// test 50% key miss
			for i := 0; i < 2*MAP_LOAD; i += 2 {
				m.Store(i, i)
			}
			return m
		},
		func(b *testing.B, m *sync.Map, drain *any) {
			for i := 0; i < b.N; i++ {
				*drain, _ = m.Load(i % (2 * MAP_LOAD))
			}
		},
	)
}

func BenchmarkLFMapRange1000000(b *testing.B) {
	benchPar(b,
		func() *Map[int, int] {
			m := New[int, int]()
			for i := 0; i < MAP_LOAD; i++ {
				m.Set(i, i)
			}
			return m
		},
		func(b *testing.B, m *Map[int, int], _ *any) {
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

func BenchmarkSyncMapRange1000000(b *testing.B) {
	benchPar(b,
		func() *sync.Map {
			m := &sync.Map{}
			for i := 0; i < MAP_LOAD; i++ {
				m.Store(i, i)
			}
			return m
		},
		func(b *testing.B, m *sync.Map, _ *any) {
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
