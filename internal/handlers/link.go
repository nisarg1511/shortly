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

const (
	statusBadRequest = "Bad Request"
	statusSuccess    = "Success!"
	statusFailed     = "Failed!"
)

func NewLink(svc *services.LinkService) *Link {
	return &Link{
		svc: svc,
	}
}

func (l *Link) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.URLShortenRequest

	var res models.APIResponse

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Status = statusBadRequest
		res.Data = "Invalid request format."
		json.NewEncoder(w).Encode(res)
		return
	}

	//URL validation
	uri, err := url.ParseRequestURI(req.URL)
	if err != nil || uri.Scheme == "" || uri.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		res.Status = statusBadRequest
		res.Data = "Invalid url format"
		json.NewEncoder(w).Encode(res)
		return
	}

	code, err := l.svc.Shorten(r.Context(), req)

	if err != nil {
		res.Status = statusFailed
		res.Data = "Could not shorten url."
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost:8000"
		log.Printf("[LINK HANDLER] No host configured, falling back to default.")
	}

	w.WriteHeader(http.StatusCreated)
	res.Status = statusSuccess
	res.Data = models.URLShortenResponse{
		ShortURL: host + "/" + code,
	}

	json.NewEncoder(w).Encode(res)
}

func (l *Link) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	w.Header().Set("Content-Type", "application/json")

	var res models.APIResponse

	if code == "" {
		res.Status = statusBadRequest
		res.Data = "Invalid request."
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	longURL, err := l.svc.GetURLFromHash(r.Context(), code)
	if err != nil {
		res.Status = statusFailed
		res.Data = "Try again later."
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusTemporaryRedirect)
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
