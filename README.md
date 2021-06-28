# Sqlc experiment

Each folder contains an experiment on [sqlc](github.com/kyleconroy/sqlc).
The experiments can be run with `go test` and a few environment variables:
`DB_HOST=localhost DB_USER=user DB_PASS=pass go test ./...`.

- [bulk](./bulk) test bulk insert and update
- [zeroes](./zeros) test zero values for various SQL types
