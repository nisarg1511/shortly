package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nisarg1511/shortly/internal/models"
)

type LinkStore struct {
	db *sql.DB
}

var store map[string]models.URLShortenRequest

func NewLinkStore(db *sql.DB) *LinkStore {
	store = make(map[string]models.URLShortenRequest)
	return &LinkStore{
		db: db,
	}
}

func (s *LinkStore) Create(ctx context.Context, link models.URLShortenRequest) error {
	if _, exists := store[link.Code]; exists {
		return errors.New("Link with code already exist!")
	}

	store[link.Code] = link
	return nil
}

func (s *LinkStore) GetByHash(ctx context.Context, hash string) (string, error) {
	entry := store[hash]

	if entry.Code == "" {
		return "", errors.New("Invalid code!")
	}
	return entry.URL, nil
}
