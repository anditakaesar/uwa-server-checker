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
	var containerID string
	var err error
	if ctx.Update.Message != nil {
		containerID, err = h.getContainerIDFromMessage(ctx.Update.Message.Text, StartContainerPrefix)
		if err != nil {
			return fmt.Errorf("couldn't get containerID from message: %v", err)
		}
	} else if ctx.Update.CallbackQuery != nil {
		containerID, err = h.getContainerIDFromMessage(ctx.Update.CallbackQuery.Data, StartContainerPrefix)
		if err != nil {
			return fmt.Errorf("couldn't get containerID from CallbackQuery: %v", err)
		}
	} else {
		return fmt.Errorf("couldn't start container with empty ID")
	}

	err = h.Docker.StartContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, fmt.Sprintf("Container %s Started", containerID), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (h *Handler) StopContainer(b *gotgbot.Bot, ctx *ext.Context) error {
	var containerID string
	var err error
	if ctx.Update.Message != nil {
		containerID, err = h.getContainerIDFromMessage(ctx.Update.Message.Text, StartContainerPrefix)
		if err != nil {
			return fmt.Errorf("couldn't get containerID from message: %v", err)
		}
	} else if ctx.Update.CallbackQuery != nil {
		containerID, err = h.getContainerIDFromMessage(ctx.Update.CallbackQuery.Data, StartContainerPrefix)
		if err != nil {
			return fmt.Errorf("couldn't get containerID from CallbackQuery: %v", err)
		}
	} else {
		return fmt.Errorf("couldn't start container with empty ID")
	}

	err = h.Docker.StopContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, fmt.Sprintf("Container %s Stoped", containerID), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to message: %w", err)
	}

	return nil
}
