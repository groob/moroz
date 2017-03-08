package moroz

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"

	"github.com/groob/moroz/santa"
)

type Endpoints struct {
	Preflight    endpoint.Endpoint
	RuleDownload endpoint.Endpoint
	EventUpload  endpoint.Endpoint
}

func MakeServerEndpoints(svc santa.Service, logger kitlog.Logger) Endpoints {
	preflightLogger := kitlog.With(logger, "method", "Preflight")
	ruleLogger := kitlog.With(logger, "method", "RuleDownload")
	eventLogger := kitlog.With(logger, "method", "EventUpload")

	preflightEndpoint := makePreflightEndpoint(svc)
	ruleDownloadEndpoint := makeRuleDownloadEndpoint(svc)
	eventUploadEndpoint := makeEventUploadEndpoint(svc)

	return Endpoints{
		Preflight:    EndpointLoggingMiddleware(preflightLogger)(preflightEndpoint),
		RuleDownload: EndpointLoggingMiddleware(ruleLogger)(ruleDownloadEndpoint),
		EventUpload:  EndpointLoggingMiddleware(eventLogger)(eventUploadEndpoint),
	}
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger kitlog.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
