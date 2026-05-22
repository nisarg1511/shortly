package store

import "database/sql"

type Store struct {
	Links *LinkStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Links: NewLinkStore(db),
		// Redis or cache stores can also be injected here later
	}
}
