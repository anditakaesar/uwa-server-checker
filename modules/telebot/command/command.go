package command

import (
	"fmt"
	"os/exec"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
)

const rejectMessage string = "You are not authorized to use this bot"

var sendMessageOpts = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
}

type Command struct {
	Docker docker.Interface
	Env    *env.Environment
}

func (cmd Command) isValidUser(userID string) bool {
	return slices.Contains(cmd.Env.ValidUserIDs(), userID)
}

func (cmd Command) Get(b *gotgbot.Bot, ctx *ext.Context) error {
	//TODO: create middleware to protect the bot
	userID := fmt.Sprint(ctx.EffectiveUser.Id)
	if cmd.isValidUser(userID) {
		//validMessage := fmt.Sprintf("Hello @%s, I'm @%s. I <b>repeat</b> all your messages.", user.Username, b.User.Username)
		_, err := ctx.EffectiveMessage.Reply(b, cmd.returnGetCommand(), &sendMessageOpts)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
		return nil
	}

	_, err := ctx.EffectiveMessage.Reply(b, rejectMessage, &sendMessageOpts)
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
	return string(out)
}

func (cmd Command) Containers(b *gotgbot.Bot, ctx *ext.Context) error {
	userID := fmt.Sprint(ctx.EffectiveUser.Id)
	if cmd.isValidUser(userID) {
		containers, err := cmd.Docker.GetContainerList()
		if err != nil {
			return fmt.Errorf("failed to get container list: %w", err)
		}
		validMessage := "<b>List of containers</b>\n"
		for _, container := range containers {
			validMessage += fmt.Sprintf("ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n\n", container.ID, container.GetNames(), container.State)
		}

		_, err = ctx.EffectiveMessage.Reply(b, validMessage, &sendMessageOpts)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
		return nil
	}

	return nil
}
