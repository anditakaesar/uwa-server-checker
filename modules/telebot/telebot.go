package telebot

import (
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot/handler"
)

type Telebot struct {
	Bot        *gotgbot.Bot
	Updater    *ext.Updater
	Dispatcher *ext.Dispatcher
	Docker     docker.Interface
	Env        *env.Environment
	Log        logger.Interface
}

func New(log logger.Interface, docker docker.Interface) (*Telebot, error) {
	env := env.New()
	bot, err := gotgbot.NewBot(env.BotToken(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("couldn't start telebot with err: %v", err), err)
		return nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Error(fmt.Sprint("an error occurred while handling update:", err.Error()), err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	return &Telebot{
		Bot:        bot,
		Updater:    updater,
		Dispatcher: dispatcher,
		Docker:     docker,
		Env:        env,
		Log:        log,
	}, nil
}

func (telebot *Telebot) InitHandlers() {
	cmd := &handler.Handler{
		Docker: telebot.Docker,
		Env:    telebot.Env,
	}

	defaultMiddlewares := []MiddlewareFunc{
		telebot.LoggingMiddleware,
		telebot.ValidUserMiddleware,
	}

	// Commands
	telebot.AddCommandHandler("get", cmd.Get, defaultMiddlewares...)
	telebot.AddCommandHandler("containers", cmd.Containers, defaultMiddlewares...)

	// Messages
	telebot.AddMessagePrefixHandler(handler.StartContainerPrefix, cmd.StartContainer, defaultMiddlewares...)
	telebot.AddMessagePrefixHandler(handler.StopContainerPrefix, cmd.StopContainer, defaultMiddlewares...)
}

func (telebot *Telebot) AddCommandHandler(command string, handlerFunc HandlerFunc, middlewareFuncs ...MiddlewareFunc) {
	var wrappedHandlerFunc handlers.Response
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		wrappedHandlerFunc = middlewareFuncs[i](handlerFunc)
	}

	telebot.Dispatcher.AddHandler(
		handlers.NewCommand(command, wrappedHandlerFunc))
}

func (telebot *Telebot) AddMessagePrefixHandler(prefix string, handlerFunc HandlerFunc, middlewareFuncs ...MiddlewareFunc) {
	var wrappedHandlerFunc handlers.Response
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		wrappedHandlerFunc = middlewareFuncs[i](handlerFunc)
	}

	telebot.Dispatcher.AddHandler(
		handlers.NewMessage(message.HasPrefix(prefix), wrappedHandlerFunc))
}

func (telebot *Telebot) Run() {
	telebot.InitHandlers()
	err := telebot.Updater.StartPolling(telebot.Bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		telebot.Log.Error(fmt.Sprintf("failed to start polling: "+err.Error()), err)
		panic(err)
	}

	telebot.Log.Info(fmt.Sprintf("%s has been started...\n", telebot.Bot.User.Username))
	telebot.Updater.Idle()
}
