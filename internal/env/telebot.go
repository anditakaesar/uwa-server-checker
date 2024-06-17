package env

import (
	"os"
	"strings"
)

var botToken *string

func readBotToken() {
	botApiToken := os.Getenv("BotToken")
	botToken = &botApiToken
}

func (e *Environment) BotToken() string {
	return *botToken
}

func readValidUserIDs() {
	envValue := os.Getenv("ValidUserIDs")
	userIDs := strings.Split(envValue, ",")
	validUserIds = userIDs
}

func (e *Environment) ValidUserIDs() []string {
	return validUserIds
}

func (e *Environment) TelebotGetCommand() string {
	return os.Getenv("TelebotGetCommand")
}
