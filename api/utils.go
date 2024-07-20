package api

import (
	"encoding/json"
	"net/http"
)

// jMap is a map that easily converts to JSON
type jMap map[string]any

func (jm jMap) Json() []byte {
	j, _ := json.Marshal(jm)
	return j
}

// jSlice is a slice that easily converts to JSON
type jSlice[T any] []T

func (js jSlice[T]) Json() []byte {
	j, _ := json.Marshal(js)
	return j
}

// all API responses are JSON
func middlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
