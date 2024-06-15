package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	fileAPPVersion    string = ".APP_VERSION"
	defaultAppVersion string = "0.0.0"
)

var appVersion *string

type Environment struct {
	appVersion string
}

func New() *Environment {
	if appVersion == nil {
		readAppVersion()
	}
	return &Environment{
		appVersion: *appVersion,
	}
}

func (e *Environment) AppName() string {
	return os.Getenv("AppName")
}

func (e *Environment) AppPort() string {
	return os.Getenv("AppPort")
}

func (e *Environment) AppVersion() string {
	return e.appVersion
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
