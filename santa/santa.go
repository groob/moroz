package santa

import (
	"context"
	"io"
)

type Service interface {
	Preflight(ctx context.Context, machineID string, p PreflightPayload) (*Preflight, error)
	RuleDownload(ctx context.Context, machineID string) ([]Rule, error)
	UploadEvent(ctx context.Context, machineID string, body io.ReadCloser) error
}

type Datastore interface {
	AllConfigs() (*ConfigCollection, error)
	Config(machineID string) (*Config, error)
}

type Config struct {
	MachineID      string `toml:"machine_id"`
	ClientMode     string `toml:"client_mode"`
	BlacklistRegex string `toml:"blacklist_regex"`
	WhitelistRegex string `toml:"whitelist_regex"`
	BatchSize      int    `toml:batch_size"`
	Rules          []Rule `toml:"rules"`
}

type ConfigCollection []*Config
