package telebot

import (
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

const rejectMessage string = "You are not authorized to use this bot"

type HandlerFunc func(b *gotgbot.Bot, ctx *ext.Context) error
type MiddlewareFunc func(handlerFunc HandlerFunc) func(b *gotgbot.Bot, ctx *ext.Context) error

var sendMessageOpts = gotgbot.SendMessageOpts{
	ParseMode: "HTML",
}

func (telebot *Telebot) isValidUser(userID int64) bool {
	userIDStr := fmt.Sprint(userID)
	return slices.Contains(telebot.Env.ValidUserIDs(), userIDStr)
}

func (telebot *Telebot) LoggingMiddleware(handlerFunc HandlerFunc) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		telebot.Log.Info(fmt.Sprintf("request from userID: %d (%s)", ctx.EffectiveUser.Id, ctx.EffectiveUser.Username), zap.Any("effectiveMessage", ctx.EffectiveMessage))
		return handlerFunc(b, ctx)
	}
}

func (telebot *Telebot) ValidUserMiddleware(handlerFunc HandlerFunc) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if telebot.isValidUser(ctx.EffectiveUser.Id) {
			return handlerFunc(b, ctx)
		}
		_, err := ctx.EffectiveMessage.Reply(b, rejectMessage, &sendMessageOpts)
		if err != nil {
			return fmt.Errorf("failed to send rejection message:  %w", err)
		}
		return nil
	}
}
