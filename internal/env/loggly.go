package env

import "os"

func LogglyBaseUrl() string {
	return os.Getenv("LogglyBaseUrl")
}

func LogglyToken() string {
	return os.Getenv("LogglyToken")
}

func LogglyTag() string {
	return os.Getenv("LogglyTag")
}
