package beelink

import (
	"fmt"
	"net/http"

	internalRouter "github.com/anditakaesar/uwa-server-checker/internal/router"
)

const modulePrefix string = "beelink"

type Handler struct{}

func New(router *internalRouter.Router) {
	handler := &Handler{}

	router.AddEndpointWithPrefix(modulePrefix, getInfo(handler))
	router.AddEndpointWithPrefix(modulePrefix, getStatus(handler))
}

func getInfo(h *Handler) internalRouter.Endpoint {
	return internalRouter.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "info",
		Handler:    h.getInfo(),
	}
}

func getStatus(h *Handler) internalRouter.Endpoint {
	return internalRouter.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "status",
		Handler:    h.getStatus(),
	}
}

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
