package initializer

import (
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"github.com/anditakaesar/uwa-server-checker/internal/logger"
	"github.com/anditakaesar/uwa-server-checker/internal/router"
	"github.com/anditakaesar/uwa-server-checker/modules/beelink"
)

type Initializer struct {
	Router *router.Router
	Log    logger.Interface
}

func New(r *router.Router) *Initializer {
	log := NewLogger()
	r.Log = log
	return &Initializer{
		Router: r,
		Log:    log,
	}
}

func NewLogger() logger.Interface {
	newClient := http.Client{}
	logglyCore := logger.NewLogglyZapCore(logger.NewLogglyLogWriter(
		logger.LogglyLogWriterDependency{
			HttpClient:    &newClient,
			BaseUrl:       env.LogglyBaseUrl(),
			CustomerToken: env.LogglyToken(),
			Tag:           env.LogglyTag(),
		},
	))

	internalLogger := logger.BuildNewLogger(logglyCore)
	return internalLogger
}

func (i *Initializer) InitModules() error {
	beelink.New(i.Router)
	return nil
}
