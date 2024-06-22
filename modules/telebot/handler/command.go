package handler

import (
	"fmt"
	"os/exec"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
)

const (
	StartContainerPrefix string = "/startContainer_"
	StopContainerPrefix  string = "/stopContainer_"
	Splitter             string = "_"

	ContainerRunningState string = "running"
	ContainerExitedState  string = "exited"
)

type Handler struct {
	Docker docker.Interface
	Env    *env.Environment
}

func (h *Handler) Get(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, h.returnGetCommand(), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func (h *Handler) returnGetCommand() string {
	cmdStr := h.Env.TelebotGetCommand()
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return fmt.Sprintf("failed to execute command: %s, with error: %+v", cmdStr, err)
	}

	if len(out) == 0 {
		return "tunnel is not running, press start and try again"
	}

	return string(out)
}

func (h *Handler) Containers(b *gotgbot.Bot, ctx *ext.Context) error {
	containers, err := h.Docker.GetContainerList()
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}
	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	}

	for _, container := range containers {
		validMessage := fmt.Sprintf("ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n\n", container.ID, container.GetNames(), container.State)
		switch container.State {
		case ContainerExitedState:
			validMessage += fmt.Sprintf("%s%s", StartContainerPrefix, container.ID)
		case ContainerRunningState:
			validMessage += fmt.Sprintf("%s%s", StopContainerPrefix, container.ID)
		}
		_, err = ctx.EffectiveMessage.Chat.SendMessage(b, validMessage, &sendMessageOpts)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
	}

	return nil
}
