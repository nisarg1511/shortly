package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request - Recieved an ill formed request.", http.StatusBadRequest)
		return
	}

	//URL validation
	uri, err := url.ParseRequestURI(req.URL)
	if err != nil || uri.Scheme == "" || uri.Host == "" {
		http.Error(w, "Invalid url format", http.StatusBadRequest)
		return
	}

	code, err := l.svc.Shorten(r.Context(), req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	var res models.APIResponse
	res.Status = "success"
	host := os.Getenv("HOST")

	if host == "" {
		host = "localhost:8000"
		log.Printf("[LINK HANDLER] No host configured, falling back to default.")
	}

	res.Data = models.URLShortenResponse{
		ShortURL: host + "/" + code,
	}

	json.NewEncoder(w).Encode(res)
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
