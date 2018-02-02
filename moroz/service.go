package moroz

import (
	"io"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

type SantaService struct {
	global      *santa.Config
	repo        santa.Datastore
	eventWriter io.Writer
}

func NewService(ds santa.Datastore, eventPath string) (*SantaService, error) {
	global, err := ds.Config("global")
	if err != nil {
		return nil, err
	}
	return &SantaService{
		global:      global,
		repo:        ds,
		eventWriter: santaEventWriter(eventPath),
	}, nil
}

func santaEventWriter(path string) io.Writer {
	events := &lumberjack.Logger{
		Filename: path,
	}
	return events
}

type Endpoints struct {
	PreflightEndpoint    endpoint.Endpoint
	RuleDownloadEndpoint endpoint.Endpoint
	EventUploadEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(svc santa.Service) Endpoints {
	return Endpoints{
		PreflightEndpoint:    makePreflightEndpoint(svc),
		RuleDownloadEndpoint: makeRuleDownloadEndpoint(svc),
		EventUploadEndpoint:  makeEventUploadEndpoint(svc),
	}
}
