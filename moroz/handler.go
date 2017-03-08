package moroz

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/moroz/santa"
)

type Handlers struct {
	Preflight    http.Handler
	RuleDownload http.Handler
	EventUpload  http.Handler
}

func MakeAPIHandler(svc santa.Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerAfter(
			kithttp.SetContentType("application/json; charset=utf-8"),
		),
	}

	endpoints := MakeServerEndpoints(svc, logger)
	h := makeHandlers(endpoints, opts)

	r := mux.NewRouter()
	r.Handle("/v1/santa/preflight/{id}", h.Preflight).Methods("POST").Name("preflight")
	r.Handle("/v1/santa/eventupload/{id}", h.EventUpload).Methods("POST").Name("eventupload")
	r.Handle("/v1/santa/ruledownload/{id}", h.RuleDownload).Methods("POST").Name("ruledownload")
	r.Handle("/v1/santa/postflight/{id}", http.HandlerFunc(nopHandler)).Methods("POST").Name("postflight")
	return r
}

func nopHandler(w http.ResponseWriter, r *http.Request) {
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	out, err := httputil.DumpRequest(r, false)
	if err != nil {
		kitlog.NewLogfmtLogger(os.Stderr).Log("err", err)
		return
	}
	fmt.Println(string(out))
}

func makeHandlers(e Endpoints, opts []kithttp.ServerOption) Handlers {
	newServer := func(e endpoint.Endpoint, decodeFn kithttp.DecodeRequestFunc) http.Handler {
		return kithttp.NewServer(e, decodeFn, encodeResponse, opts...)
	}
	return Handlers{
		Preflight:    newServer(e.Preflight, decodePreflightRequest),
		RuleDownload: newServer(e.RuleDownload, decodeRuleRequest),
		EventUpload:  newServer(e.EventUpload, decodeEventUpload),
	}
}
