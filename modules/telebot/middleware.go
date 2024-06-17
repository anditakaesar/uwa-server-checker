package telebot

import (
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const rejectMessage string = "You are not authorized to use this bot"

var sendMessageOpts = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
}

func (telebot *Telebot) isValidUser(userID int64) bool {
	userIDStr := fmt.Sprint(userID)
	return slices.Contains(telebot.Env.ValidUserIDs(), userIDStr)
}

func (telebot *Telebot) ValidUserOnly(handler func(b *gotgbot.Bot, ctx *ext.Context) error) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if telebot.isValidUser(ctx.EffectiveUser.Id) {
			return handler(b, ctx)
		}
		_, err := ctx.EffectiveMessage.Reply(b, rejectMessage, &sendMessageOpts)
		if err != nil {
			return fmt.Errorf("failed to send rejection message:  %w", err)
		}
		return nil
	}
}
