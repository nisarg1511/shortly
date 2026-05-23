package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/nisarg1511/shortly/internal/models"
)

type LinkStore struct {
	db *sql.DB
}

func NewLinkStore(db *sql.DB) *LinkStore {
	return &LinkStore{
		db: db,
	}
}

func (s *LinkStore) Create(ctx context.Context, link models.URLShortenRequest) error {
	query := `INSERT INTO urls(short_code,original_url) VALUES ($1,$2)`
	_, err := s.db.Exec(query, link.Code, link.URL)
	if err != nil {
		log.Printf("[LINK STORE] Last link insertion failed. Reason:%v\n", err)
		return err
	}
	return nil
}

func (s *LinkStore) GetByHash(ctx context.Context, hash string) (string, error) {
	var url string
	query := `SELECT original_url FROM urls WHERE short_code = $1`
	err := s.db.QueryRow(query, hash).Scan(&url)
	if err != nil {
		log.Printf("[LINK STORE] Couldn't find url for given code. Reason:%v\n", err)
		return "", err
	}
	return url, nil
}
