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
	MachineID                     string `toml:"machine_id"`
	ClientMode                    string `toml:"client_mode"`
	BlacklistRegex                string `toml:"blacklist_regex"`
	WhitelistRegex                string `toml:"whitelist_regex"`
	BatchSize                     int    `toml:batch_size"`
	Rules                         []Rule `toml:"rules"`
	EnableBundles                 bool   `toml:"enable_bundles"`
	EnabledTransitiveWhitelisting bool   `toml:"enabled_transitive_whitelisting"`
}

type Rule struct {
	RuleType      string `json:"rule_type" toml:"rule_type"`
	Policy        string `json:"policy" toml:"policy"`
	SHA256        string `json:"sha256" toml:"sha256"`
	CustomMessage string `json:"custom_msg,omitempty" toml:"custom_msg"`
}

type ConfigCollection []*Config
