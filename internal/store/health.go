package store

import (
	"context"
	"database/sql"
	"log"
)

type HealthStore struct {
	db *sql.DB
}

func NewHealthStore(db *sql.DB) *HealthStore {
	return &HealthStore{
		db: db,
	}
}

func (s *HealthStore) CheckDatabaseHealth(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		log.Printf("[HEALTH STORE] CRITICAL: Could not ping database. Reason:%v\n", err)
		return err
	}
	return nil
}
