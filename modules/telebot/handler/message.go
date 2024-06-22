package handler

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) getContainerIDFromMessage(messageStr string, prefix string) (string, error) {
	values := strings.Split(messageStr, Splitter)
	if len(values) < 1 {
		return "", fmt.Errorf("couldn't get message values of %s with prefix %s", messageStr, prefix)
	}

	return values[1], nil
}

func (h *Handler) StartContainer(b *gotgbot.Bot, ctx *ext.Context) error {
	containerID, err := h.getContainerIDFromMessage(ctx.Update.Message.Text, StartContainerPrefix)
	if err != nil {
		return err
	}

	err = h.Docker.StartContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Container Started", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (h *Handler) StopContainer(b *gotgbot.Bot, ctx *ext.Context) error {
	containerID, err := h.getContainerIDFromMessage(ctx.Update.Message.Text, StopContainerPrefix)
	if err != nil {
		return err
	}

	err = h.Docker.StopContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Container Stopped", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to message: %w", err)
	}

	return nil
}
