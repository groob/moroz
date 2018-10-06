package santa

import (
	"encoding/json"
	"fmt"
)

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
