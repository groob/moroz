package santa

import (
	"encoding/json"
	"fmt"
)

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

type PreflightPayload struct {
	OSBuild              string `json:"os_build"`
	SantaVersion         string `json:"santa_version"`
	Hostname             string `json:"hostname"`
	OSVersion            string `json:"os_version"`
	CertificateRuleCount int    `json:"certificate_rule_count"`
	BinaryRuleCount      int    `json:"binary_rule_count"`
	ClientMode           string `json:"client_mode"`
	SerialNumber         string `json:"serial_number"`
	PrimaryUser          string `json:"primary_user"`
}

type ClientMode int

const (
	Monitor ClientMode = iota
	Lockdown
)

func (c ClientMode) MarshalJSON() ([]byte, error) {
	var mode string
	switch c {
	case Monitor:
		mode = "MONITOR"
	case Lockdown:
		mode = "LOCKDOWN"
	default:
		return nil, fmt.Errorf("unknown client_mode %d", c)
	}
	return json.Marshal(mode)
}

type Preflight struct {
	WhitelistRegex                string     `json:"whitelist_regex"`
	BlacklistRegex                string     `json:"blacklist_regex"`
	BatchSize                     int        `json:"batch_size"`
	ClientMode                    ClientMode `json:"client_mode"`
	EnableBundles                 bool       `json:"enable_bundles"`
	EnabledTransitiveWhitelisting bool       `json:"enabled_transitive_whitelisting"`
}
