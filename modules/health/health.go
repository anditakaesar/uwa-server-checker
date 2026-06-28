package health

import (
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/internal/router"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot"
)

const modulePrefix string = "health"

func NewModule(router *router.Router, telebot *telebot.Telebot) {
	handler := &Handler{
		Telebot: telebot,
	}

	router.AddEndpointWithPrefix(modulePrefix, getInfo(handler))
	router.AddEndpointWithPrefix(modulePrefix, getStatus(handler))
	router.AddEndpointWithPrefix("api", sendMessage(handler))
}

func getInfo(h *Handler) router.Endpoint {
	return router.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "info",
		Handler:    h.getInfo(),
	}
}

func getStatus(h *Handler) router.Endpoint {
	return router.Endpoint{
		HTTPMethod: http.MethodGet,
		Path:       "status",
		Handler:    h.getStatus(),
	}
}

func sendMessage(h *Handler) router.Endpoint {
	return router.Endpoint{
		HTTPMethod: http.MethodPost,
		Path:       "message",
		Handler:    h.sendMessage(),
	}
}
