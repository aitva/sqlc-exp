package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	exp "github.com/aitva/sqlc-exp"
	"github.com/aitva/sqlc-exp/zeros/query"
	_ "github.com/lib/pq"
)

var queries *query.Queries

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

	queries, err = query.Prepare(context.Background(), db)
	if err != nil {
		fmt.Printf("fail to prepare queries: %v\n", err)
		return 1
	}

	err = createAuthors()
	if err != nil {
		fmt.Printf("fail to create authors: %v\n", err)
		return 1
	}

	return m.Run()
}

func createAuthors() error {
	authors := []query.CreateAuthorParams{
		{
			ID:    0,
			Name:  "Octavia E. Butler",
			Birth: time.Date(1947, 6, 24, 0, 0, 0, 0, time.UTC),
			Bio: toNullString(`Octavia Estelle Butler was an American ` +
				`science fiction author. A multiple recipient of both the Hugo ` +
				`and Nebula awards, she became in 1995 the first science-fiction ` +
				`writer to receive a MacArthur Fellowship.`),
		},
		{
			ID:    1,
			Name:  "Arthur C. Clarke",
			Birth: time.Date(1917, 12, 16, 0, 0, 0, 0, time.UTC),
			Bio: toNullString(`Sir Arthur Charles Clarke was an English ` +
				`science-fiction writer, science writer, futurist, inventor, ` +
				`undersea explorer, and television series host.`),
		},
		{
			ID:    2,
			Name:  "Kae Tempest",
			Birth: time.Date(1985, 12, 22, 0, 0, 0, 0, time.UTC),
			Bio: toNullString(`Kae Tempest is an English spoken word ` +
				`performer, poet, recording artist, novelist and playwright. ` +
				`At the age of 16, they were accepted into the exclusive ` +
				`BRIT School for Performing Arts and Technology in Croydon.`),
		},
	}
	for i := range authors {
		a, err := queries.CreateAuthor(context.Background(), authors[i])
		if err != nil {
			return err
		}
		authors[i].ID = a.ID
	}
	return nil
}

func TestListAuthors(t *testing.T) {
	tests := []struct {
		params  query.ListAuthorsParams
		expects []string
	}{
		{
			params: query.ListAuthorsParams{
				Name: "Octavia E. Butler",
			},
			expects: []string{"Octavia E. Butler"},
		},
		{
			params: query.ListAuthorsParams{
				Birth: time.Date(1917, 12, 16, 0, 0, 0, 0, time.UTC),
			},
			expects: []string{"Arthur C. Clarke"},
		},
		{
			params: query.ListAuthorsParams{
				Bio: "BRIT School",
			},
			expects: []string{"Kae Tempest"},
		},
	}
	for i, tt := range tests {
		authors, err := queries.ListAuthors(context.Background(), tt.params)
		if err != nil {
			t.Fatalf("tests[%d]: fail to list authors: %v", i, err)
		}

		if len(tt.expects) != len(authors) {
			t.Fatalf("tests[%d]: missing authors; expects=%d got=%d", i, len(tt.expects), len(authors))
		}

		expects := tt.expects
		for _, a := range authors {
			var ok bool
			expects, ok = filterAuthor(expects, &a)
			if !ok {
				t.Fatalf("tests[%d]: unexpected user; got=%v expects=%v", i, a.Name, tt.expects)
			}
		}
	}
}

func filterAuthor(expects []string, a *query.Author) ([]string, bool) {
	for i, e := range expects {
		if e == a.Name {
			return append(expects[:i], expects[i+1:]...), true
		}
	}
	return expects, false
}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
