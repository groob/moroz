package santa

import (
	"github.com/pkg/errors"
)

type Config struct {
	MachineID string `toml:"machine_id,omitempty"`
	Preflight
	Rules []Rule `toml:"rules"`
}

type Rule struct {
	RuleType      RuleType `json:"rule_type" toml:"rule_type"`
	Policy        Policy   `json:"policy" toml:"policy"`
	SHA256        string   `json:"sha256" toml:"sha256"`
	CustomMessage string   `json:"custom_msg,omitempty" toml:"custom_msg,omitempty"`
}

type Preflight struct {
	ClientMode                    ClientMode `json:"client_mode" toml:"client_mode"`
	BlacklistRegex                string     `json:"blacklist_regex" toml:"blacklist_regex"`
	WhitelistRegex                string     `json:"whitelist_regex" toml:"whitelist_regex"`
	BatchSize                     int        `json:"batch_size" toml:"batch_size"`
	EnableBundles                 bool       `json:"enable_bundles" toml:"enable_bundles"`
	EnabledTransitiveWhitelisting bool       `json:"enabled_transitive_whitelisting" toml:"enabled_transitive_whitelisting"`
}

type PreflightPayload struct {
	OSBuild              string     `json:"os_build"`
	SantaVersion         string     `json:"santa_version"`
	Hostname             string     `json:"hostname"`
	OSVersion            string     `json:"os_version"`
	CertificateRuleCount int        `json:"certificate_rule_count"`
	BinaryRuleCount      int        `json:"binary_rule_count"`
	ClientMode           ClientMode `json:"client_mode"`
	SerialNumber         string     `json:"serial_number"`
	PrimaryUser          string     `json:"primary_user"`
}

type RuleType int

const (
	Binary RuleType = iota
	Certificate
)

func (r *RuleType) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BINARY":
		*r = Binary
	case "CERTIFICATE":
		*r = Certificate
	default:
		return errors.Errorf("unknown rule_type value %q", t)
	}
	return nil
}

func (r RuleType) MarshalText() ([]byte, error) {
	switch r {
	case Binary:
		return []byte("BINARY"), nil
	case Certificate:
		return []byte("CERTIFICATE"), nil
	default:
		return nil, errors.Errorf("unknown rule_type %d", r)
	}
}

type Policy int

const (
	Blacklist Policy = iota
	Whitelist
	WhitelistCompiler
)

func (p *Policy) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BLACKLIST":
		*p = Blacklist
	case "WHITELIST":
		*p = Whitelist
	case "WHITELIST_COMPILER":
		*p = WhitelistCompiler
	default:
		return errors.Errorf("unknown policy value %q", t)
	}
	return nil
}

func (p Policy) MarshalText() ([]byte, error) {
	switch p {
	case Blacklist:
		return []byte("BLACKLIST"), nil
	case Whitelist:
		return []byte("WHITELIST"), nil
	case WhitelistCompiler:
		return []byte("WHITELIST_COMPILER"), nil
	default:
		return nil, errors.Errorf("unknown policy %d", p)
	}
}

type ClientMode int

const (
	Monitor ClientMode = iota
	Lockdown
)

func (c *ClientMode) UnmarshalText(text []byte) error {
	switch mode := string(text); mode {
	case "MONITOR":
		*c = Monitor
	case "LOCKDOWN":
		*c = Lockdown
	default:
		return errors.Errorf("unknown client_mode value %q", mode)
	}
	return nil
}

func (c ClientMode) MarshalText() ([]byte, error) {
	switch c {
	case Monitor:
		return []byte("MONITOR"), nil
	case Lockdown:
		return []byte("LOCKDOWN"), nil
	default:
		return nil, errors.Errorf("unknown client_mode %d", c)
	}
}
