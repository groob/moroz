package moroz

import (
	"github.com/go-kit/kit/log"
	"github.com/groob/moroz/santa"
)

type Middleware func(santa.Service) santa.Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next santa.Service) santa.Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	next   santa.Service
}
