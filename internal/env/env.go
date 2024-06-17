package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	fileAPPVersion    string = ".release-please-manifest.json"
	defaultAppVersion string = "0.0.0"
)

var (
	envPtr       *Environment
	appVersion   *string
	apiToken     *string
	validUserIds []string
)

type Environment struct{}

func New() *Environment {
	if envPtr != nil {
		return envPtr
	}
	readAppVersion()
	readApiToken()
	readBotToken()
	readValidUserIDs()

	envPtr = &Environment{}
	return envPtr
}

func (e *Environment) AppName() string {
	return os.Getenv("AppName")
}

func (e *Environment) AppPort() string {
	return os.Getenv("AppPort")
}

func (e *Environment) Env() string {
	return os.Getenv("Env")
}

func (e *Environment) GetAddrPort() string {
	return fmt.Sprint(":", e.AppPort())
}

func (e *Environment) AppVersion() string {
	return *appVersion
}

func readAppVersionErr(err error) {
	if err != nil {
		fmt.Printf("error while reading app version: %v", err)
	}

	defaultVer := defaultAppVersion
	appVersion = &defaultVer
}

func readAppVersion() {
	jsonFile, err := os.Open(fileAPPVersion)
	if err != nil {
		readAppVersionErr(err)
	}
	defer jsonFile.Close()

	byteValues, err := io.ReadAll(jsonFile)
	if err != nil {
		readAppVersionErr(err)
	}

	version := map[string]string{}

	err = json.Unmarshal(byteValues, &version)
	if err != nil {
		readAppVersionErr(err)
	}

	versionStr, ok := version["."]
	if !ok {
		readAppVersionErr(errors.New("failed to fetch data"))
	}

	if strings.TrimSpace(versionStr) == "" {
		defaultVer := defaultAppVersion
		appVersion = &defaultVer
	}
	appVersion = &versionStr
}

func (e *Environment) ApiToken() string {
	return *apiToken
}

func readApiToken() {
	envApiToken := os.Getenv("ApiToken")
	if envApiToken == "" {
		apiTokenStr := uuid.New().String()
		apiToken = &apiTokenStr
	}

	apiToken = &envApiToken
}
