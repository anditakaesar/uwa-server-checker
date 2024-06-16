package telebot

import (
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/modules/telebot/command"
)

type Telebot struct {
	Bot        *gotgbot.Bot
	Updater    *ext.Updater
	Dispatcher *ext.Dispatcher
	Log        logger.Interface
}

func New(log logger.Interface) (*Telebot, error) {
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
		Log:        log,
	}, nil
}

func (telebot *Telebot) InitCommands() {
	command := command.Command{}
	telebot.Dispatcher.AddHandler(handlers.NewCommand("start", command.Start))
}

func (telebot *Telebot) Run() {
	telebot.InitCommands()
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
