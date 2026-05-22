package models

type URLShortenRequest struct {
	URL    string `json:"url"`
	Code   string `json:"code"`
	Expiry int64  `json:"expiry"`
}
