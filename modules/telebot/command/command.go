package command

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
)

const (
	StartCallbackPrefix string = "startCallback_"
	StopCallbackPrefix  string = "stopCallback_"
	Splitter            string = "_"
)

var sendMessageOpts = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
}

type Command struct {
	Docker docker.Interface
	Env    *env.Environment
}

func (cmd Command) Get(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, cmd.returnGetCommand(), &sendMessageOpts)
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func (cmd Command) returnGetCommand() string {
	cmdStr := cmd.Env.TelebotGetCommand()
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return fmt.Sprintf("failed to execute command: %s, with error: %+v", cmdStr, err)
	}

	if len(out) == 0 {
		return "tunnel is not running, press start and try again"
	}

	return string(out)
}

func (cmd Command) getContainerIDFromCallback(callbackStr string, prefix string) (string, error) {
	values := strings.Split(callbackStr, Splitter)
	if len(values) < 1 {
		return "", fmt.Errorf("couldn't get callback values of %s with prefix %s", callbackStr, prefix)
	}

	return values[1], nil
}

func (cmd Command) GetStartCB(b *gotgbot.Bot, ctx *ext.Context) error {
	sendMessageOptsCopy := sendMessageOpts
	containerID, err := cmd.getContainerIDFromCallback(ctx.Update.CallbackQuery.Data, StartCallbackPrefix)
	if err != nil {
		return err
	}

	err = cmd.Docker.StartContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, "Container Started", &sendMessageOptsCopy)
	if err != nil {
		return fmt.Errorf("failed to send callback message: %w", err)
	}

	return nil
}

func (cmd Command) GetStopCB(b *gotgbot.Bot, ctx *ext.Context) error {
	sendMessageOptsCopy := sendMessageOpts
	containerID, err := cmd.getContainerIDFromCallback(ctx.Update.CallbackQuery.Data, StopCallbackPrefix)
	if err != nil {
		return err
	}

	err = cmd.Docker.StopContainer(containerID)
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, "Container Stopped", &sendMessageOptsCopy)
	if err != nil {
		return fmt.Errorf("failed to send callback message: %w", err)
	}

	return nil
}

func (cmd Command) Containers(b *gotgbot.Bot, ctx *ext.Context) error {
	containers, err := cmd.Docker.GetContainerList()
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}
	for _, container := range containers {
		sendMessageOptsCopy := sendMessageOpts
		sendMessageOptsCopy.ReplyMarkup = gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Start",
						CallbackData: fmt.Sprintf("%s%s", StartCallbackPrefix, container.ID),
					},
					{
						Text:         "Stop",
						CallbackData: fmt.Sprintf("%s%s", StopCallbackPrefix, container.ID),
					},
				},
			},
		}
		validMessage := fmt.Sprintf("ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n\n", container.ID, container.GetNames(), container.State)
		_, err = ctx.EffectiveMessage.Reply(b, validMessage, &sendMessageOptsCopy)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
	}

	return nil
}
