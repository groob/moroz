package moroz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func AddHTTPRoutes(r *mux.Router, e Endpoints, logger log.Logger) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(
			httptransport.SetContentType("application/json; charset=utf-8"),
		),
	}

	// POST     /v1/santa/preflight/:id			preflight request.
	// POST     /v1/santa/ruledownload/:id		request rule updates.
	// POST     /v1/santa/eventupload/:id		upload event.
	// POST     /v1/santa/postflight/:id		postflight request.

	r.Methods("POST").Path("/v1/santa/preflight/{id}").Handler(httptransport.NewServer(
		e.PreflightEndpoint,
		decodePreflightRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/santa/ruledownload/{id}").Handler(httptransport.NewServer(
		e.RuleDownloadEndpoint,
		decodeRuleRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/santa/eventupload/{id}").Handler(httptransport.NewServer(
		e.EventUploadEndpoint,
		decodeEventUpload,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/santa/postflight/{id}").Handler(httptransport.NewServer(
		e.PostflightEndpoint,
		decodePostflightRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// errBadRoute is used for mux errors
var errBadRoute = errors.New("bad route")

func machineIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return "", errBadRoute
	}
	return id, nil
}

type failer interface {
	Failed() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	errMap := map[string]interface{}{"error": err.Error()}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}

	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc.Encode(errMap)
}
