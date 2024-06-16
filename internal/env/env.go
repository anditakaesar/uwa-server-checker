package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	fileAPPVersion    string = ".APP_VERSION"
	defaultAppVersion string = "0.0.0"
)

var appVersion *string
var apiToken *string

type Environment struct{}

func New() *Environment {
	readAppVersion()
	readApiToken()

	return &Environment{}
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
	file, err := os.Open(fileAPPVersion)
	if err != nil {
		readAppVersionErr(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		readAppVersionErr(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		readAppVersionErr(err)
	}

	appVersionString := string(buffer)
	if strings.TrimSpace(appVersionString) == "" {
		defaultVer := defaultAppVersion
		appVersion = &defaultVer
	}
	appVersion = &appVersionString
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
