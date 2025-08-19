package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/aitva/sqlc-exp/bulk/query"
)

//go:generate sqlc generate

type Bulk struct {
	queries *query.Queries
	counter atomic.Int64
}

func New(db *sql.DB) (*Bulk, error) {
	queries, err := query.Prepare(context.Background(), db)
	if err != nil {
		return nil, fmt.Errorf("prepare: %v", err)
	}
	return &Bulk{
		queries: queries,
	}, nil
}

func (b *Bulk) ID() int64    { return b.counter.Add(1) }
func (b *Bulk) MaxID() int64 { return b.counter.Load() }

func (b *Bulk) CreateAuthor(ctx context.Context, a Author) error {
	return b.queries.CreateAuthor(ctx, query.CreateAuthorParams{
		ID:   a.ID,
		Name: a.Name,
		Bio:  ptrToNullString(a.Bio),
	})
}

func (b *Bulk) BulkCreateAuthor(ctx context.Context, authors []Author) error {
	tmps := query.CreateAuthorsParams{
		IDs:   make([]int64, len(authors)),
		Names: make([]string, len(authors)),
		Bios:  make([]sql.NullString, len(authors)),
	}

	for i, a := range authors {
		tmps.IDs[i] = a.ID
		tmps.Names[i] = a.Name
		tmps.Bios[i] = ptrToNullString(a.Bio)
	}

	return b.queries.CreateAuthors(ctx, tmps)
}

func (b *Bulk) UpdateAuthor(ctx context.Context, a Author) error {
	return b.queries.UpdateAuthor(ctx, query.UpdateAuthorParams{
		ID:   a.ID,
		Name: a.Name,
		Bio:  ptrToNullString(a.Bio),
	})
}

func (b *Bulk) BulkUpdateAuthor(ctx context.Context, authors []Author) error {
	tmps := query.UpdateAuthorsParams{
		IDs:   make([]int64, len(authors)),
		Names: make([]string, len(authors)),
		Bios:  make([]sql.NullString, len(authors)),
	}

	for i, a := range authors {
		tmps.IDs[i] = a.ID
		tmps.Names[i] = a.Name
		tmps.Bios[i] = ptrToNullString(a.Bio)
	}

	return b.queries.UpdateAuthors(ctx, tmps)
}

func ptrToNullString(s *string) sql.NullString {
	var tmp sql.NullString
	if s != nil {
		tmp = sql.NullString{
			Valid:  true,
			String: *s,
		}
	}
	return tmp
}

func ptr[T any](t T) *T { return &t }
