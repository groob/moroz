package santa

type Rule struct {
	RuleType      string `json:"rule_type" toml:"rule_type"`
	Policy        string `json:"policy" toml:"policy"`
	SHA256        string `json:"sha256" toml:"sha256"`
	CustomMessage string `json:"custom_msg,omitempty" toml:"custom_msg"`
}
