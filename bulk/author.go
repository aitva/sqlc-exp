package bulk

import "fmt"

type Author struct {
	ID   int64
	Name string
	Bio  *string
}

// NewAuthor creates a new [Author] with the given ID.
func NewAuthor(id int64) Author {
	return Author{
		ID:   id,
		Name: fmt.Sprintf("Author %d", id),
		Bio:  ptr(fmt.Sprintf("Author %d is quite good.", id)),
	}
}

// UpdateAuthor updates an existing [Author].
func UpdateAuthor(id int64) Author {
	return Author{
		ID:   id,
		Name: fmt.Sprintf("Author %d", id),
		Bio:  ptr(fmt.Sprintf("Author %d is excellent!", id)),
	}
}
