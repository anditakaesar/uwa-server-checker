package health

import (
	"fmt"
	"net/http"
)

type Handler struct{}

func (handler *Handler) getInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"message":"success"}`)
	})
}

func (handler *Handler) getStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"message":"success", "path":"status"}`)
	})
}
