package batch

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	exp "github.com/aitva/sqlc-exp"
	"github.com/aitva/sqlc-exp/batch/query"
	_ "github.com/lib/pq"
)

var db *exp.DB
var counter int64

// TestMain setup and teardown the database before the tests.
func TestMain(m *testing.M) {
	var err error
	db, err = exp.LoadDB()
	if err != nil {
		fmt.Printf("fail to load db: %v\n", err)
		os.Exit(1)
	}

	err = db.Up()
	if err != nil {
		fmt.Printf("fail to create db: %v\n", err)
	}

	code := m.Run()

	err = db.Drop()
	if err != nil {
		fmt.Printf("fail to drop db: %v\n", err)
	}
	if code == 0 && err != nil {
		code = 1
	}

	os.Exit(code)
}

func BenchmarkCreateAuthor(b *testing.B) {
	q := query.New(db)
	for i := 0; i < b.N; i++ {
		id := atomic.AddInt64(&counter, 1)
		err := q.CreateAuthor(context.Background(), query.CreateAuthorParams{
			ID:   id,
			Name: fmt.Sprintf("Author %d", id),
			Bio: sql.NullString{
				Valid:  true,
				String: fmt.Sprintf("Author %d is an exceptional person.", id),
			},
		})
		if err != nil {
			b.Fatalf("fail to create author: %v", err)
		}
	}
}

func BenchmarkCreateAuthors10(b *testing.B) { benchmarkCreateAuthors(b, 10) }

func BenchmarkCreateAuthors100(b *testing.B) { benchmarkCreateAuthors(b, 100) }

func BenchmarkCreateAuthors1000(b *testing.B) { benchmarkCreateAuthors(b, 1000) }

func benchmarkCreateAuthors(b *testing.B, count int) {
	var params query.CreateAuthorsParams
	q := query.New(db)
	for i := 0; i < b.N; i++ {
		// We aggregate the columns into arrays.
		id := atomic.AddInt64(&counter, 1)
		params.Ids = append(params.Ids, id)
		params.Names = append(params.Names, fmt.Sprintf("Author %d", id))
		params.Bios = append(params.Bios, fmt.Sprintf("Author %d is an exceptional person.", id))
		if len(params.Ids) < count {
			continue
		}

		// We create multiple records at once.
		err := q.CreateAuthors(context.Background(), params)
		if err != nil {
			b.Fatalf("fail to create authors: %v", err)
		}
		params.Ids = params.Ids[:0]
		params.Names = params.Names[:0]
		params.Bios = params.Bios[:0]
	}
}

func BenchmarkUpdateAuthor(b *testing.B) {
	q := query.New(db)
	for i := 0; i < b.N; i++ {
		id := int64(i) % counter
		err := q.UpdateAuthor(context.Background(), query.UpdateAuthorParams{
			ID:   id,
			Name: fmt.Sprintf("Author %d updated", id),
			Bio: sql.NullString{
				Valid:  true,
				String: fmt.Sprintf("Author %d is still an exceptional person.", id),
			},
		})
		if err != nil {
			b.Fatalf("fail to update author: %v", err)
		}
	}
}

func BenchmarkUpdateAuthors10(b *testing.B) { benchmarkUpdateAuthors(b, 10) }

func BenchmarkUpdateAuthors100(b *testing.B) { benchmarkUpdateAuthors(b, 100) }

func BenchmarkUpdateAuthors1000(b *testing.B) { benchmarkUpdateAuthors(b, 1000) }

func benchmarkUpdateAuthors(b *testing.B, count int) {
	var params query.UpdateAuthorsParams
	q := query.New(db)
	for i := 0; i < b.N; i++ {
		// We aggregate the columns into arrays.
		id := int64(i) % counter
		params.Ids = append(params.Ids, id)
		params.Names = append(params.Names, fmt.Sprintf("Author %d updated", id))
		params.Bios = append(params.Bios, fmt.Sprintf("Author %d is still an exceptional person.", id))
		if len(params.Ids) < count {
			continue
		}

		// We update multiple records at once.
		err := q.UpdateAuthors(context.Background(), params)
		if err != nil {
			b.Fatalf("fail to create authors: %v", err)
		}
		params.Ids = params.Ids[:0]
		params.Names = params.Names[:0]
		params.Bios = params.Bios[:0]
	}
}
