package env

import "os"

var botToken *string

func readBotToken() {
	botApiToken := os.Getenv("BotToken")
	botToken = &botApiToken
}

func (e *Environment) BotToken() string {
	return *botToken
}
