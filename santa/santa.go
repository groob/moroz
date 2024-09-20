// Package santa defines types for a Santa sync server.
package santa

import (
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
// https://github.com/google/santa/blob/ff0efe952b2456b52fad2a40e6eedb0931e6bdf7/docs/development/sync-protocol.md#rules-objects
type Rule struct {
	RuleType      RuleType `json:"rule_type" toml:"rule_type"`
	Policy        Policy   `json:"policy" toml:"policy"`
	Identifier    string   `json:"identifier" toml:"identifier"`
	CustomMessage string   `json:"custom_msg,omitempty" toml:"custom_msg,omitempty"`
	CustomUrl     string   `json:"custom_url,omitempty" toml:"custom_url,omitempty"`
	// TODO: add support for the following fields
	// CreationTime          float64  `json:"creation_time,omitempty" toml:"creation_time,omitempty"`
	// FileBundleBinaryCount int      `json:"file_bundle_binary_count,omitempty" toml:"file_bundle_binary_count,omitempty"`
	// FileBundleHash        string   `json:"file_bundle_hash,omitempty" toml:"file_bundle_hash,omitempty"`
}

// Preflight represents sync response sent to a Santa client by the sync server.
// https://github.com/google/santa/blob/344a35aaf63c24a56f7a021ce18ecab090584da3/docs/development/sync-protocol.md#preflight-response
type Preflight struct {
	ClientMode            ClientMode `json:"client_mode" toml:"client_mode"`
	BlockedPathRegex      string     `json:"blocked_path_regex" toml:"blocked_path_regex"`
	AllowedPathRegex      string     `json:"allowed_path_regex" toml:"allowed_path_regex"`
	BatchSize             int        `json:"batch_size" toml:"batch_size"`
	EnableAllEventUpload  bool       `json:"enable_all_event_upload" toml:"enable_all_event_upload"`
	EnableBundles         bool       `json:"enable_bundles" toml:"enable_bundles"`
	EnableTransitiveRules bool       `json:"enable_transitive_rules" toml:"enable_transitive_rules"`
	CleanSync             bool       `json:"clean_sync" toml:"clean_sync,omitempty"`
	FullSyncInterval      int        `json:"full_sync_interval" toml:"full_sync_interval"`
	// TODO: add support for sync_type and deprecate clean_sync
	//	SyncType                 string     `json:"sync_type" toml:"sync_type,omitempty"`
	// TODO: add in support for the following fields
	//	BlockUSBMount            bool   `json:"block_usb_mount" toml:"block_usb_mount,omitempty"`
	//	RemountUSBMode           string `json:"remount_usb_mode" toml:"remount_usb_mode,omitempty"`
	//	OverrideFileAccessAction string `json:"override_file_access_action" toml:"override_file_access_action,omitempty"`
}

// A PreflightPayload represents the request sent by a santa client to the sync server.
// https://github.com/google/santa/blob/344a35aaf63c24a56f7a021ce18ecab090584da3/docs/development/sync-protocol.md#preflight-request
type PreflightPayload struct {
	SerialNumber         string     `json:"serial_num"`
	Hostname             string     `json:"hostname"`
	OSVersion            string     `json:"os_version"`
	OSBuild              string     `json:"os_build"`
	ModelIdentifier      string     `json:"model_identifier"`
	SantaVersion         string     `json:"santa_version"`
	PrimaryUser          string     `json:"primary_user"`
	BinaryRuleCount      int        `json:"binary_rule_count"`
	CertificateRuleCount int        `json:"certificate_rule_count"`
	CompilerRuleCount    int        `json:"compiler_rule_count"`
	TransitiveRuleCount  int        `json:"transitive_rule_count"`
	TeamIDRuleCount      int        `json:"teamid_rule_count"`
	SigningIDRuleCount   int        `json:"signingid_rule_count"`
	CdHashRuleCount      int        `json:"cdhash_rule_count"`
	ClientMode           ClientMode `json:"client_mode"`
	RequestCleanSync     bool       `json:"request_clean_sync"`
}

// Postflight represents sync response sent to a Santa client by the sync server.
// Currently, this is a no-op.
type Postflight struct {
	NoOp struct{}
}

// A PostflightPayload represents the request sent by a santa client to the sync server.
// https://github.com/google/santa/blob/344a35aaf63c24a56f7a021ce18ecab090584da3/docs/development/sync-protocol.md#postflight-request
type PostflightPayload struct {
	RulesReceived  int `json:"rules_received"`
	RulesProcessed int `json:"rules_processed"`
}

// EventPayload represents derived metadata for events uploaded with the UploadEvent endpoint.
type EventPayload struct {
	FileSHA   string  `json:"file_sha256"`
	UnixTime  float64 `json:"execution_time"`
	EventInfo EventUploadEvent
}

// EventUploadRequest encapsulation of an /eventupload POST body sent by a Santa client
type EventUploadRequest struct {
	Events []EventUploadEvent `json:"events"`
}

// EventUploadEvent is a single event entry
// https://github.com/google/santa/blob/344a35aaf63c24a56f7a021ce18ecab090584da3/docs/development/sync-protocol.md#event-objects
type EventUploadEvent struct {
	CurrentSessions              []string       `json:"current_sessions"`
	Decision                     string         `json:"decision"`
	ExecutingUser                string         `json:"executing_user"`
	ExecutionTime                float64        `json:"execution_time"`
	FileBundleBinaryCount        int64          `json:"file_bundle_binary_count"`
	FileBundleExecutableRelPath  string         `json:"file_bundle_executable_rel_path"`
	FileBundleHash               string         `json:"file_bundle_hash"`
	FileBundleHashMilliseconds   float64        `json:"file_bundle_hash_millis"`
	FileBundleID                 string         `json:"file_bundle_id"`
	FileBundleName               string         `json:"file_bundle_name"`
	FileBundlePath               string         `json:"file_bundle_path"`
	FileBundleShortVersionString string         `json:"file_bundle_version_string"`
	FileBundleVersion            string         `json:"file_bundle_version"`
	FileName                     string         `json:"file_name"`
	FilePath                     string         `json:"file_path"`
	FileSHA256                   string         `json:"file_sha256"`
	LoggedInUsers                []string       `json:"logged_in_users"`
	ParentName                   string         `json:"parent_name"`
	ParentProcessID              int            `json:"ppid"`
	ProcessID                    int            `json:"pid"`
	QuarantineAgentBundleID      string         `json:"quarantine_agent_bundle_id"`
	QuarantineDataUrl            string         `json:"quarantine_data_url"`
	QuarantineRefererUrl         string         `json:"quarantine_referer_url"`
	QuarantineTimestamp          float64        `json:"quarantine_timestamp"`
	SigningChain                 []SigningEntry `json:"signing_chain"`
	SigningID                    string         `json:"signing_id"`
	TeamID                       string         `json:"team_id"`
	CdHash                       string         `json:"cd_hash"`
}

// SigningEntry is optionally present when an event includes a binary that is signed
type SigningEntry struct {
	CertificateName    string `json:"cn"`
	Organization       string `json:"org"`
	OrganizationalUnit string `json:"ou"`
	SHA256             string `json:"sha256"`
	ValidFrom          int    `json:"valid_from"`
	ValidUntil         int    `json:"valid_until"`
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

	// TeamID rules are the 10-character identifier issued by Apple and tied to developer accounts/organizations.
	// This is an even more powerful rule with broader reach than individual certificate rules.
	// ie. EQHXZ8M8AV for Google
	TeamID

	// Signing IDs are arbitrary identifiers under developer control that are given to a binary at signing time.
	// Because the signing IDs are arbitrary, the Santa rule identifier must be prefixed with the Team ID associated
	// with the Apple developer certificate used to sign the application.
	// ie. EQHXZ8M8AV:com.google.Chrome
	SigningID

	// CDHash rules use a binary's code directory hash as an identifier. This is the most specific rule in Santa.
	// The code directory hash identifies a specific version of a program, similar to a file hash.
	// Note that the operating system evaluates the cdhash lazily, only verifying pages of code when they're mapped in.
	// This means that it is possible for a file hash to change, but a binary could still execute as long as modified
	// pages are not mapped in. Santa only considers CDHash rules for processes that have CS_KILL or CS_HARD
	// codesigning flags set to ensure that a process will be killed if the CDHash was tampered with
	// (assuming the system has SIP enabled).
	CdHash
)

func (r *RuleType) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BINARY":
		*r = Binary
	case "CERTIFICATE":
		*r = Certificate
	case "TEAMID":
		*r = TeamID
	case "SIGNINGID":
		*r = SigningID
	case "CDHASH":
		*r = CdHash
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
	case SigningID:
		return []byte("SIGNINGID"), nil
	case CdHash:
		return []byte("CDHASH"), nil
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
	Remove
)

func (p *Policy) UnmarshalText(text []byte) error {
	switch t := string(text); t {
	case "BLOCKLIST":
		*p = Blocklist
	case "ALLOWLIST":
		*p = Allowlist
	case "ALLOWLIST_COMPILER":
		*p = AllowlistCompiler
	case "REMOVE":
		*p = Remove
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
	case Remove:
		return []byte("REMOVE"), nil
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
