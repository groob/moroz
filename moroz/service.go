package moroz

import (
	"context"
	"io"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/groob/moroz/santa"
)

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

type SantaService struct {
	global      *santa.Config
	repo        santa.Datastore
	eventWriter io.Writer
}

func (svc *SantaService) Preflight(ctx context.Context, machineID string, p santa.PreflightPayload) (*santa.Preflight, error) {
	config := svc.config(machineID)
	pre := &santa.Preflight{
		ClientMode:     toClientMode(config.ClientMode),
		BlacklistRegex: config.BlacklistRegex,
		WhitelistRegex: config.WhitelistRegex,
		BatchSize:      config.BatchSize,
	}
	return pre, nil
}

func (svc *SantaService) RuleDownload(ctx context.Context, machineID string) ([]santa.Rule, error) {
	config := svc.config(machineID)
	return config.Rules, nil
}

func (svc *SantaService) UploadEvent(ctx context.Context, machineID string, body io.ReadCloser) error {
	_, err := io.Copy(svc.eventWriter, body)
	defer body.Close()
	return err
}

func (svc *SantaService) config(machineID string) *santa.Config {
	var config *santa.Config
	var err error
	config, err = svc.repo.Config(machineID)
	if err != nil {
		config = svc.global
	}
	return config
}

func toClientMode(from string) santa.ClientMode {
	switch from {
	case "MONITOR":
		return santa.Monitor
	case "LOCKDOWN":
		return santa.Lockdown
	default:
		return santa.Monitor
	}

}
