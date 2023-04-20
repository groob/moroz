// Package santa defines types for a Santa sync server.
package santa

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Config represents the combination of the Preflight configuration and Rules
// for a given MachineID.
type Config struct {
	MachineID string `toml:"machine_id,omitempty"`
	Preflight
	Rules []Rule `toml:"rules"`
}

// Rule is a Santa rule.
// Full documentation: https://github.com/google/santa/blob/01df4623c7c534568ca3d310129455ff71cc3eef/Docs/details/rules.md
type Rule struct {
	RuleType      RuleType `json:"rule_type" toml:"rule_type"`
	Policy        Policy   `json:"policy" toml:"policy"`
	SHA256        string   `json:"sha256,omitempty" toml:"sha256,omitempty"`
	Identifier    string   `json:"identifier,omitempty" toml:"identifier,omitempty"`
	CustomMessage string   `json:"custom_msg,omitempty" toml:"custom_msg,omitempty"`
}

// Preflight representssync response sent to a Santa client by the sync server.
type Preflight struct {
	ClientMode                    ClientMode `json:"client_mode" toml:"client_mode"`
	BlocklistRegex                string     `json:"blocklist_regex" toml:"blocklist_regex"`
	AllowlistRegex                string     `json:"allowlist_regex" toml:"allowlist_regex"`
	BatchSize                     int        `json:"batch_size" toml:"batch_size"`
	EnableBundles                 bool       `json:"enable_bundles" toml:"enable_bundles"`
	EnabledTransitiveAllowlisting bool       `json:"enabled_transitive_allowlisting" toml:"enabled_transitive_allowlisting"`
}

// A PreflightPayload represents the request sent by a santa client to the sync server.
type PreflightPayload struct {
	OSBuild              string     `json:"os_build"`
	SantaVersion         string     `json:"santa_version"`
	Hostname             string     `json:"hostname"`
	OSVersion            string     `json:"os_version"`
	CertificateRuleCount int        `json:"certificate_rule_count"`
	BinaryRuleCount      int        `json:"binary_rule_count"`
	ClientMode           ClientMode `json:"client_mode"`
	SerialNumber         string     `json:"serial_num"`
	PrimaryUser          string     `json:"primary_user"`
}

// EventPayload represents derived metadata for events uploaded with the UploadEvent endpoint.
type EventPayload struct {
	FileSHA  string          `json:"file_sha256"`
	UnixTime float64         `json:"execution_time"`
	Content  json.RawMessage `json:"-"`
}

// RuleType represents a Santa rule type.
type RuleType int

const (
	// Binary rules use the SHA-256 hash of the entire binary as an identifier.
	Binary RuleType = iota

	// Certificate rules are formed from the SHA-256 fingerprint of an X.509 leaf signing certificate.
	// This is a powerful rule type that has a much broader reach than an individual binary rule .
	// A signing certificate can sign any number of binaries.
	Certificate

	TeamID
)

func (r *RuleType) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BINARY":
		*r = Binary
	case "CERTIFICATE":
		*r = Certificate
	case "TEAMID":
		*r = TeamID
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
	case TeamID:
		return []byte("TEAMID"), nil
	default:
		return nil, errors.Errorf("unknown rule_type %d", r)
	}
}

// Policy represents the Santa Rule Policy.
type Policy int

const (
	Blocklist Policy = iota
	Allowlist

	// AllowlistCompiler is a Transitive allowlist policy which allows allowlisting binaries created by
	// a specific compiler. EnabledTransitiveAllowlisting must be set to true in the Preflight first.
	AllowlistCompiler
)

func (p *Policy) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BLOCKLIST":
		*p = Blocklist
	case "ALLOWLIST":
		*p = Allowlist
	case "ALLOWLIST_COMPILER":
		*p = AllowlistCompiler
	default:
		return errors.Errorf("unknown policy value %q", t)
	}
	return nil
}

func (p Policy) MarshalText() ([]byte, error) {
	switch p {
	case Blocklist:
		return []byte("BLOCKLIST"), nil
	case Allowlist:
		return []byte("ALLOWLIST"), nil
	case AllowlistCompiler:
		return []byte("ALLOWLIST_COMPILER"), nil
	default:
		return nil, errors.Errorf("unknown policy %d", p)
	}
}

// ClientMode specifies which mode the Santa client will evaluate rules in.
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
