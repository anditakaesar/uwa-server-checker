package telebot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/anditakaesar/uwa-server-checker/adapter/docker"
	"github.com/anditakaesar/uwa-server-checker/internal"
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
	telebot.AddCommandHandler("get", cmd.Get, defaultMiddlewares...)
	telebot.AddCommandHandler("containers", cmd.InitializeReplyContainerPaging, defaultMiddlewares...)

	// Callbacks
	telebot.AddCallbackHandler(handler.ContainerPagingPrefix, cmd.ProcessCallbackContainerPaging, defaultMiddlewares...)

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

func (telebot *Telebot) run() error {
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

func (telebot *Telebot) shutdown() error {
	_, err := telebot.Bot.Close(nil)
	return err
}

type Module struct{}

func (mod *Module) Start(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error, dep internal.Dependency) {
	log := logger.GetLogInstance()
	wg.Add(1)
	go func() {
		defer wg.Done()

		botObj, err := New(Dependency{
			Docker: dep.Adapter.Docker,
		})
		if err != nil {
			log.Error("failed to initialize bot object", err)
			errCh <- fmt.Errorf("initialize bot server error: %w", err)
			return
		}

		// Start server in separate goroutine
		go func() {
			if err := botObj.run(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("TelebotModule start error: %w", err)
			}
		}()

		// Wait for shutdown signal
		<-ctx.Done()
		log.Info("Shutting down TelebotModule...")
		if err := botObj.shutdown(); err != nil {
			errCh <- fmt.Errorf("TelebotModule shutdown error: %w", err)
		}
	}()
}
