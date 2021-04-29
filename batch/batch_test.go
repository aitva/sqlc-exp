package batch

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	exp "github.com/aitva/sqlc-exp"
	_ "github.com/lib/pq"
)

var db *exp.DB

func TestMain(m *testing.M) {
	var err error
	db, err = exp.LoadDB()
	if err != nil {
		fmt.Printf("fail to load db: %v", err)
		os.Exit(1)
	}

}

func testMain(m *testing.M) int {
	return m.Run()
}

func BenchmarkCreateAuthor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}

func BenchmarkCreateAuthors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}
