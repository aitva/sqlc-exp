package bulk

import (
	"fmt"
	"os"
	"testing"

	exp "github.com/aitva/sqlc-exp"
	_ "github.com/lib/pq"
)

var bulk *Bulk

// TestMain setup and teardown the database before the tests.
func TestMain(m *testing.M) {
	db, err := exp.LoadDB()
	if err != nil {
		fmt.Printf("fail to load db: %v\n", err)
		os.Exit(1)
	}

	code := testMain(db, m)

	err = db.Drop()
	if err != nil {
		fmt.Printf("fail to drop db: %v\n", err)
	}
	if code == 0 && err != nil {
		code = 1
	}

	os.Exit(code)
}

func testMain(db *exp.DB, m *testing.M) int {
	err := db.Up()
	if err != nil {
		fmt.Printf("fail to create db: %v\n", err)
		return 1
	}

	bulk, err = New(db.DB)
	if err != nil {
		fmt.Printf("fail to create bulk: %v\n", err)
		return 1
	}

	return m.Run()
}

func BenchmarkCreateAuthor(b *testing.B) {
	for b.Loop() {
		err := bulk.CreateAuthor(b.Context(), NewAuthor(bulk.ID()))
		if err != nil {
			b.Fatalf("fail to create author: %v", err)
		}
	}

	// Compute time per operation by multiplying b.N by the lenght of
	// the dataset.
	//us := b.Elapsed().Microseconds() / int64(b.N)
	//b.ReportMetric(float64(us), "μs/op")
	ns := b.Elapsed().Nanoseconds() / int64(b.N)
	b.ReportMetric(float64(ns), "ns/op")
}

func BenchmarkCreateAuthors10(b *testing.B) { benchmarkCreateAuthors(b, 10) }

func BenchmarkCreateAuthors100(b *testing.B) { benchmarkCreateAuthors(b, 100) }

func BenchmarkCreateAuthors1000(b *testing.B) { benchmarkCreateAuthors(b, 1000) }

func BenchmarkCreateAuthors10000(b *testing.B) { benchmarkCreateAuthors(b, 10000) }

func benchmarkCreateAuthors(b *testing.B, count int) {
	for b.Loop() {
		authors := make([]Author, count)
		for i := range authors {
			authors[i] = NewAuthor(bulk.ID())
		}

		err := bulk.BulkCreateAuthor(b.Context(), authors)
		if err != nil {
			b.Fatalf("fail to bulk create authors: %v", err)
		}
	}

	// Compute time per operation by multiplying b.N by the lenght of
	// the dataset.
	//us := b.Elapsed().Microseconds() / int64(b.N*count)
	//b.ReportMetric(float64(us), "μs/op")
	ns := b.Elapsed().Nanoseconds() / int64(b.N*count)
	b.ReportMetric(float64(ns), "ns/op")
}

func BenchmarkUpdateAuthor(b *testing.B) {
	max := bulk.MaxID()
	id := int64(0)
	for b.Loop() {
		err := bulk.UpdateAuthor(b.Context(), UpdateAuthor(id))
		if err != nil {
			b.Fatalf("fail to update author: %v", err)
		}
		id = (id + 1) % max
	}

	// Compute time per operation by multiplying b.N by the lenght of
	// the dataset.
	//us := b.Elapsed().Microseconds() / int64(b.N)
	//b.ReportMetric(float64(us), "μs/op")
	ns := b.Elapsed().Nanoseconds() / int64(b.N)
	b.ReportMetric(float64(ns), "ns/op")
}

func BenchmarkUpdateAuthors10(b *testing.B) { benchmarkUpdateAuthors(b, 10) }

func BenchmarkUpdateAuthors100(b *testing.B) { benchmarkUpdateAuthors(b, 100) }

func BenchmarkUpdateAuthors1000(b *testing.B) { benchmarkUpdateAuthors(b, 1000) }

func BenchmarkUpdateAuthors10000(b *testing.B) { benchmarkUpdateAuthors(b, 10000) }

func benchmarkUpdateAuthors(b *testing.B, count int) {
	max := bulk.MaxID()
	id := int64(0)
	for b.Loop() {
		authors := make([]Author, count)
		for i := range authors {
			authors[i] = UpdateAuthor(id)
		}

		err := bulk.BulkUpdateAuthor(b.Context(), authors)
		if err != nil {
			b.Fatalf("fail to bulk update authors: %v", err)
		}

		id = (id + 1) % max
	}

	// Compute time per operation by multiplying b.N by the lenght of
	// the dataset.
	//us := b.Elapsed().Microseconds() / int64(b.N*count)
	//b.ReportMetric(float64(us), "μs/op")
	ns := b.Elapsed().Nanoseconds() / int64(b.N*count)
	b.ReportMetric(float64(ns), "ns/op")
}
