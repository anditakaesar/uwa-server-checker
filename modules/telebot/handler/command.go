package handler

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
)

const (
	StartContainerPrefix  string = "/startContainer_"
	StopContainerPrefix   string = "/stopContainer_"
	FindContainerPrefix   string = "/findcontainer"
	ContainerPagingPrefix string = "containerspg_"
	Splitter              string = "_"

	ContainerRunningState string = "running"
	ContainerExitedState  string = "exited"
)

type Handler struct {
	Docker docker.Interface
	Env    *env.Environment
}

func (h *Handler) GetOpenAddress(b *gotgbot.Bot, ctx *ext.Context) error {
	containers, err := h.Docker.GetContainersByName(h.Env.TunnelContainerName())
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}

	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	}

	row1 := []gotgbot.InlineKeyboardButton{}
	bodyMessage := strings.Builder{}
	for _, container := range containers {
		fmt.Fprintf(&bodyMessage, "ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n",
			container.ID, container.GetNames(), container.State)
		switch container.State {
		case ContainerExitedState:
			row1 = append(row1, gotgbot.InlineKeyboardButton{
				Text:         "Start",
				CallbackData: fmt.Sprintf("%s%s", StartContainerPrefix, container.ID),
			})
		case ContainerRunningState:
			row1 = append(row1, gotgbot.InlineKeyboardButton{
				Text:         "Stop",
				CallbackData: fmt.Sprintf("%s%s", StopContainerPrefix, container.ID),
			})
		}
		bodyMessage.WriteString("\n\n")
	}

	bodyMessage.WriteString(h.returnGetCommand())
	if len(row1) > 0 {
		sendMessageOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				row1,
			},
		}
	}

	_, err = ctx.EffectiveMessage.Reply(b, bodyMessage.String(), &sendMessageOpts)
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

func generatePagingCallback(prefix string, size, page int) string {
	return fmt.Sprintf("%s%d_%d", prefix, size, page)
}

func parseSizePage(data string) (int, int) {
	parsed := strings.Split(data, "_")
	if len(parsed) < 3 {
		return 1, 1
	}

	size, err := strconv.Atoi(parsed[1])
	if err != nil {
		return 1, 1
	}

	page, err := strconv.Atoi(parsed[2])
	if err != nil {
		return 1, 1
	}

	return size, page
}

func (h *Handler) FindContainerByName(b *gotgbot.Bot, ctx *ext.Context) error {
	text := strings.Split(ctx.Update.Message.Text, " ")
	if len(text) < 2 {
		return fmt.Errorf("please specify name criteria")
	}
	containers, err := h.Docker.GetContainersByName(text[1])
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}

	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	}
	bodyMessage := strings.Builder{}
	for _, container := range containers {
		fmt.Fprintf(&bodyMessage, "ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n",
			container.ID, container.GetNames(), container.State)
		switch container.State {
		case ContainerExitedState:
			fmt.Fprintf(&bodyMessage, "%s%s", StartContainerPrefix, container.ID)
		case ContainerRunningState:
			fmt.Fprintf(&bodyMessage, "%s%s", StopContainerPrefix, container.ID)
		}
		bodyMessage.WriteString("\n\n")
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, bodyMessage.String(), &sendMessageOpts)
	if err != nil {
		return fmt.Errorf("failed to send FindContainerByName message: %w", err)
	}

	return nil
}

func (h *Handler) InitializeReplyContainerPaging(b *gotgbot.Bot, ctx *ext.Context) error {
	const initSize int = 5
	nextCallback := generatePagingCallback(ContainerPagingPrefix, initSize, 2)
	sendMessageOpts := gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Next",
						CallbackData: nextCallback,
					},
				},
			},
		},
	}

	result, err := h.Docker.GetContainerListWithPaging(initSize, 1)
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}

	bodyMessage := strings.Builder{}
	for _, container := range result.List {
		fmt.Fprintf(&bodyMessage, "ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n",
			container.ID, container.GetNames(), container.State)
		switch container.State {
		case ContainerExitedState:
			fmt.Fprintf(&bodyMessage, "%s%s", StartContainerPrefix, container.ID)
		case ContainerRunningState:
			fmt.Fprintf(&bodyMessage, "%s%s", StopContainerPrefix, container.ID)
		}
		bodyMessage.WriteString("\n\n")
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, bodyMessage.String(), &sendMessageOpts)
	if err != nil {
		return fmt.Errorf("failed to send InitializeReplyContainerPaging message: %w", err)
	}

	return nil
}

func (h *Handler) ProcessCallbackContainerPaging(b *gotgbot.Bot, ctx *ext.Context) error {
	size, page := parseSizePage(ctx.Update.CallbackQuery.Data)
	result, err := h.Docker.GetContainerListWithPaging(size, page)
	if err != nil {
		return fmt.Errorf("failed to get container list: %w", err)
	}

	prevCallback := generatePagingCallback(ContainerPagingPrefix, size, page-1)
	nextCallback := generatePagingCallback(ContainerPagingPrefix, size, page+1)

	// generate Prev | Next
	row1 := []gotgbot.InlineKeyboardButton{}
	if result.HasPrev {
		row1 = append(row1, gotgbot.InlineKeyboardButton{
			Text:         "Prev",
			CallbackData: prevCallback,
		})
	}

	if result.HasNext {
		row1 = append(row1, gotgbot.InlineKeyboardButton{
			Text:         "Next",
			CallbackData: nextCallback,
		})
	}

	editMessageTextOpts := gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{ // use the previous data
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				row1,
			},
		},
	}

	bodyMessage := strings.Builder{}
	for _, container := range result.List {
		fmt.Fprintf(&bodyMessage, "ID: <b>%s</b> \nName: <b>%s</b> \nState: <i>%s</i>\n",
			container.ID, container.GetNames(), container.State)
		switch container.State {
		case ContainerExitedState:
			fmt.Fprintf(&bodyMessage, "%s%s", StartContainerPrefix, container.ID)
		case ContainerRunningState:
			fmt.Fprintf(&bodyMessage, "%s%s", StopContainerPrefix, container.ID)
		}
		bodyMessage.WriteString("\n\n")
	}

	_, _, err = ctx.EffectiveMessage.EditText(b, bodyMessage.String(), &editMessageTextOpts)
	if err != nil {
		return fmt.Errorf("failed to Edit message: %w", err)
	}

	return nil
}
