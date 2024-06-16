package beelink

import (
	"net/http"

	internalRouter "github.com/anditakaesar/uwa-server-checker/internal/router"
)

const modulePrefix string = "beelink"

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
