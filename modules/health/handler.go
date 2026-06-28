package health

import (
	"encoding/json"
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/modules/telebot"
)

type Handler struct {
	Telebot *telebot.Telebot
}

func (handler *Handler) getInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "success",
		})
	})
}

func (handler *Handler) getStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "success",
		})
	})
}

type CustomMessage struct {
	Message string `json:"message"`
}

func (handler *Handler) sendMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data CustomMessage
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Telebot.SendMessage(data.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// send actual message
		json.NewEncoder(w).Encode(data)
	})
}
