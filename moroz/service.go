package moroz

import (
	"context"
	"io"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/moroz/santa"
)

type ConfigStore interface {
	AllConfigs() ([]santa.Config, error)
	Config(machineID string) (santa.Config, error)
}

type SantaService struct {
	global      santa.Config
	repo        ConfigStore
	eventWriter io.Writer
}

func NewService(ds ConfigStore, eventPath string) (*SantaService, error) {
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

type Service interface {
	Preflight(ctx context.Context, machineID string, p santa.PreflightPayload) (*santa.Preflight, error)
	RuleDownload(ctx context.Context, machineID string) ([]santa.Rule, error)
	UploadEvent(ctx context.Context, machineID string, body io.ReadCloser) error
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

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		PreflightEndpoint:    makePreflightEndpoint(svc),
		RuleDownloadEndpoint: makeRuleDownloadEndpoint(svc),
		EventUploadEndpoint:  makeEventUploadEndpoint(svc),
	}
}
