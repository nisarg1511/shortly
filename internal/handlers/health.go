package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nisarg1511/shortly/internal/models"
	"github.com/nisarg1511/shortly/internal/services"
)

type Health struct {
	svc *services.HealthService
}

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

func NewHealth(svc *services.HealthService) *Health {
	return &Health{
		svc: svc,
	}
}

func (h *Health) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	var res models.APIResponse

	w.Header().Set("Content-Type", "application/json")
	if err := h.svc.CheckHealth(ctx); err != nil {
		res.Status = StatusUnhealthy
		res.Data = models.HealthCheckResponse{
			Database: "disconnected",
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Status = StatusHealthy
	res.Data = models.HealthCheckResponse{
		Database: "connected",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
