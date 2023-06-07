package santa

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestConfigMarshalUnmarshal(t *testing.T) {
	conf := testConfig(t, "testdata/config_a_toml.golden", (os.Getenv("REPLACE_GOLDEN") == "TRUE"))

	if have, want := conf.ClientMode, Lockdown; have != want {
		t.Errorf("have client_mode %d, want %d\n", have, want)
	}

	if have, want := conf.CleanSync, true; have != want {
		t.Errorf("have clean_sync %t, want %t\n", have, want)
	}

	if have, want := conf.FullSyncInterval, 600; have != want {
		t.Errorf("have full_sync_interval %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[0].Identifier, "2dc104631939b4bdf5d6bccab76e166e37fe5e1605340cf68dab919df58b8eda"; have != want {
		t.Errorf("have identifier %s, want %s\n", have, want)
	}

	if have, want := conf.Rules[0].RuleType, Binary; have != want {
		t.Errorf("have rule_type %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[1].RuleType, Certificate; have != want {
		t.Errorf("have rule_type %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[2].RuleType, TeamID; have != want {
		t.Errorf("have rule_tpe %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[3].RuleType, SigningID; have != want {
		t.Errorf("have rule_type %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[0].Policy, Blocklist; have != want {
		t.Errorf("have policy %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[1].Policy, Allowlist; have != want {
		t.Errorf("have policy %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[4].Policy, AllowlistCompiler; have != want {
		t.Errorf("have policy %d, want %d\n", have, want)
	}

	if have, want := conf.Rules[5].Policy, Remove; have != want {
		t.Errorf("have policy %d, want %d\n", have, want)
	}
}

func testConfig(t *testing.T, path string, replace bool) Config {
	t.Helper()

	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("loading config from path %q, err = %q\n", path, err)
	}

	var conf Config
	if err := toml.Unmarshal(file, &conf); err != nil {
		t.Fatalf("unmarshal config from path %q, err = %q\n", path, err)
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(&conf); err != nil {
		t.Fatalf("encode config from path %q, err = %q\n", path, err)
	}

	if replace {
		if err := ioutil.WriteFile(path, buf.Bytes(), os.ModePerm); err != nil {
			t.Fatalf("replace config at path %q, err = %q\n", path, err)
		}
		return testConfig(t, path, false)
	}

	if !bytes.Equal(file, buf.Bytes()) {
		t.Errorf("marshaling config to %q failed\nEXPECTED:\n%s\nGOT:\n%s\n", path, string(file), buf.Bytes())

	}

	return conf
}
