package telebot

import (
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
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
}

type Dependency struct {
	Docker docker.Interface
}

func New(dep Dependency) (*Telebot, error) {
	env := env.New()
	log := logger.GetLogInstance()
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
		Docker:     dep.Docker,
		Env:        env,
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
	telebot.AddCommandHandler("get", cmd.GetOpenAddress, defaultMiddlewares...)
	telebot.AddCommandHandler("containers", cmd.InitializeReplyContainerPaging, defaultMiddlewares...)
	telebot.AddCommandHandler("backupbw", cmd.RunBitwardenBackup, defaultMiddlewares...)

	// Callbacks
	telebot.AddCallbackHandler(handler.ContainerPagingPrefix, cmd.ProcessCallbackContainerPaging, defaultMiddlewares...)
	telebot.AddCallbackHandler(handler.StartContainerPrefix, cmd.StartContainer, defaultMiddlewares...)
	telebot.AddCallbackHandler(handler.StopContainerPrefix, cmd.StopContainer, defaultMiddlewares...)

	// Messages
	telebot.AddMessagePrefixHandler(handler.StartContainerPrefix, cmd.StartContainer, defaultMiddlewares...)
	telebot.AddMessagePrefixHandler(handler.StopContainerPrefix, cmd.StopContainer, defaultMiddlewares...)
	telebot.AddMessagePrefixHandler(handler.FindContainerPrefix, cmd.FindContainerByName, defaultMiddlewares...)
}

func (telebot *Telebot) AddCommandHandler(
	command string,
	handlerFunc func(b *gotgbot.Bot, ctx *ext.Context) error,
	middlewareFuncs ...MiddlewareFunc,
) {
	wrappedHandlerFunc := handlerFunc
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		wrappedHandlerFunc = middlewareFuncs[i](wrappedHandlerFunc)
	}

	telebot.Dispatcher.AddHandler(
		handlers.NewCommand(command, wrappedHandlerFunc))
}

func (telebot *Telebot) AddCallbackHandler(
	prefix string,
	handlerFunc func(b *gotgbot.Bot, ctx *ext.Context) error,
	middlewareFuncs ...MiddlewareFunc,
) {
	wrappedHandlerFunc := handlerFunc
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		wrappedHandlerFunc = middlewareFuncs[i](wrappedHandlerFunc)
	}

	telebot.Dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(prefix), handlerFunc))
}

func (telebot *Telebot) AddMessagePrefixHandler(
	prefix string,
	handlerFunc func(b *gotgbot.Bot, ctx *ext.Context) error,
	middlewareFuncs ...MiddlewareFunc,
) {
	wrappedHandlerFunc := handlerFunc
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		wrappedHandlerFunc = middlewareFuncs[i](wrappedHandlerFunc)
	}

	telebot.Dispatcher.AddHandler(
		handlers.NewMessage(message.HasPrefix(prefix), wrappedHandlerFunc))
}

func (telebot *Telebot) Run() error {
	telebot.InitHandlers()
	log := logger.GetLogInstance()
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
		log.Error(fmt.Sprintf("failed to start polling: %s", err.Error()), err)
		return err
	}

	log.Info(fmt.Sprintf("%s has been started...\n", telebot.Bot.User.Username))
	telebot.Updater.Idle()
	return nil
}

func (telebot *Telebot) Shutdown() error {
	_, err := telebot.Bot.Close(nil)
	return err
}

func (telebot *Telebot) SendMessage(message string) error {
	if telebot.Env.NotificationChatID() == nil {
		return fmt.Errorf("ChatID unset")
	}

	_, err := telebot.Bot.SendMessage(*telebot.Env.NotificationChatID(), message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return err
	}

	return nil
}
