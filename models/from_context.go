package models

import (
	"net/http"
)

// Context keys
type contextKey string

const (
	StoreKey = contextKey("store")
)

func FromContextStore(r *http.Request) (*Store, bool) {
	store, ok := r.Context().Value(StoreKey).(*Store)
	return store, ok
}
