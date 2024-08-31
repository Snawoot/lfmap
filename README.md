# lfmap

[![Go Reference](https://pkg.go.dev/badge/github.com/Snawoot/lfmap.svg)](https://pkg.go.dev/github.com/Snawoot/lfmap)

Generic concurrent lock-free map for Golang.

Key features:

* `range` iteration over consistent snapshot of map without locking other threads and without copying all data in map to maintain that snapshot.
* Convenient [transactions](https://pkg.go.dev/github.com/Snawoot/lfmap#Map.Transaction) where you can read-write multiple keys, check some conditions and make a change based on your logic. Compared to lfmap, `sync.Map` allows only CAS operations against single key and that's it.

## Usage

See [godoc examples](https://pkg.go.dev/github.com/Snawoot/lfmap#pkg-examples).

## Benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/Snawoot/lfmap
cpu: Intel(R) N100
BenchmarkLFMapSet-4              	  165698	     11190 ns/op
BenchmarkSyncMapSet-4            	  857448	      2302 ns/op
BenchmarkLFMapGet-4              	 3221922	       365.1 ns/op
BenchmarkSyncMapGet-4            	 5302554	       189.9 ns/op
BenchmarkLFMapRange1000000-4     	       8	 142530076 ns/op
BenchmarkSyncMapRange1000000-4   	       7	 150277709 ns/op
PASS
ok  	github.com/Snawoot/lfmap	31.614s
```

So far lfmap is 2-6 times slower than `sync.Map`, mostly because of underlying immutable ds performance. However, nice transactional properties may make it useful.
