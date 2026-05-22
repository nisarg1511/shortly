package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nisarg1511/shortly/internal/models"
	"github.com/nisarg1511/shortly/internal/services"
)

type Link struct {
	svc *services.LinkService
}

func NewLink(svc *services.LinkService) *Link {
	return &Link{
		svc: svc,
	}
}

func (l *Link) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.URLShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad Request - Recieved an ill formed request.", http.StatusBadRequest)
	}
	if req.URL == "" {
		http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
		return
	}

	code, err := l.svc.Shorten(r.Context(), req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte("New Url: localhost:8000/" + code))
}

func (l *Link) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "missing 'code' parameter", http.StatusBadRequest)
		return // Stop execution
	}

	longURL, err := l.svc.GetURLFromHash(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // Stop execution
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
