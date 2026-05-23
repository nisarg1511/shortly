package services

import "github.com/nisarg1511/shortly/internal/store"

type Service struct {
	Links         *LinkService
	HealthService *HealthService
}

func NewService(storage *store.Store) *Service {
	return &Service{
		Links:         NewLinkService(storage.Links),
		HealthService: NewHealthService(storage.Health),
	}
}
