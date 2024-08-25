package lfmap

import (
	"fmt"
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
