package services

import (
	"context"

	"github.com/nisarg1511/shortly/internal/models"
	"github.com/nisarg1511/shortly/internal/store"
)

type LinkService struct {
	linkStore *store.LinkStore
}

var Counter int

func NewLinkService(ls *store.LinkStore) *LinkService {
	return &LinkService{
		linkStore: ls,
	}
}

func (s *LinkService) Shorten(ctx context.Context, shortenReq models.URLShortenRequest) (string, error) {

	if shortenReq.Code == "" {
		code := getUniqueCode(Counter)
		shortenReq.Code = code
	}

	err := s.linkStore.Create(ctx, shortenReq)
	if err != nil {
		return "", err
	}
	Counter++

	return shortenReq.Code, nil
}

func (s *LinkService) GetURLFromHash(ctx context.Context, hash string) (string, error) {
	url, err := s.linkStore.GetByHash(ctx, hash)
	if err != nil {
		return "", err
	}

	return url, nil
}

func getUniqueCode(num int) string {
	// Fix 1: Handle zero explicitly
	if num == 0 {
		return "0"
	}

	lookup := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base62 := ""

	for num > 0 {
		base62 = string(lookup[num%62]) + base62
		num /= 62
	}

	return base62
}
