package services

import (
	"context"
	"log"

	"github.com/nisarg1511/shortly/internal/store"
)

type HealthService struct {
	healthStore *store.HealthStore
}

func NewHealthService(hs *store.HealthStore) *HealthService {
	return &HealthService{
		healthStore: hs,
	}
}

func (s *HealthService) CheckHealth(ctx context.Context) error {
	if err := s.healthStore.CheckDatabaseHealth(ctx); err != nil {
		log.Printf("[HEALTH SERVICE] CRITICAL: Could not ping to database.\n")
		return err
	}
	return nil
}
