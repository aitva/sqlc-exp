# Sqlc experiment

Each folder contains an experiment on [sqlc](github.com/kyleconroy/sqlc).
The experiments can be run with `go test` and a few environment variables:
`DB_HOST=localhost DB_USER=sqlc-exp DB_PASS=sqlc-exp DB_SSLMODE=disable go test -bench . ./...`.

## Bulk

The package `bulk` measure bulk performance using Postgres arrays and `sqlc`.
The benchmark hint that using bulk insert and update is ~200 times faster than
using normal operations.

```
$ DB_HOST=localhost DB_USER=sqlc-exp DB_PASS=sqlc-exp DB_SSLMODE=disable go test -bench . ./bulk
goos: linux
goarch: amd64
pkg: github.com/aitva/sqlc-exp/bulk
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkCreateAuthor-16          	    534	  2336591 ns/op
BenchmarkCreateAuthors10-16       	    523	   239403 ns/op
BenchmarkCreateAuthors100-16      	    378	    31597 ns/op
BenchmarkCreateAuthors1000-16     	    240	     4696 ns/op
BenchmarkCreateAuthors10000-16    	     61	     1996 ns/op
BenchmarkUpdateAuthor-16          	    496	  2533515 ns/op
BenchmarkUpdateAuthors10-16       	    500	   258094 ns/op
BenchmarkUpdateAuthors100-16      	    380	    31103 ns/op
BenchmarkUpdateAuthors1000-16     	    214	     5039 ns/op
BenchmarkUpdateAuthors10000-16    	     79	     1925 ns/op
PASS
ok  	github.com/aitva/sqlc-exp/bulk	12.425s
```

## Zeroes

The package `zeroes` test a method to write queries with conditional parameters
and `sqlc`.
