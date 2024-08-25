package lfmap_test

import (
	"fmt"
	"sync"

	"github.com/Snawoot/lfmap"
)

func Example() {
	m := lfmap.New[string, int]()
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
		fmt.Printf("key = %s, value = %d\n", k, v)
	}
}
