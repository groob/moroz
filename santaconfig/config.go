package santaconfig

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/groob/moroz/santa"
	"github.com/pkg/errors"
)

func NewFileRepo(path string) *FileRepo {
	repo := FileRepo{
		configIndex: make(map[string]santa.Config),
		configPath:  path,
	}
	return &repo
}

type FileRepo struct {
	mtx         sync.RWMutex
	configIndex map[string]santa.Config
	configPath  string
}

func (f *FileRepo) updateIndex(configs []santa.Config) {
	f.configIndex = make(map[string]santa.Config, len(configs))
	for _, conf := range configs {
		f.configIndex[conf.MachineID] = conf
	}
}

func (f *FileRepo) AllConfigs(ctx context.Context) ([]santa.Config, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	configs, err := loadConfigs(f.configPath)
	if err != nil {
		return nil, err
	}
	f.updateIndex(configs)
	return configs, nil
}

func (f *FileRepo) Config(ctx context.Context, machineID string) (santa.Config, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	var conf santa.Config
	configs, err := loadConfigs(f.configPath)
	if err != nil {
		return conf, errors.Wrapf(err, "loading config for machineID %q", machineID)
	}
	f.updateIndex(configs)
	conf, ok := f.configIndex[machineID]
	if !ok {
		return conf, errors.Errorf("configuration %q not found", machineID)
	}
	return conf, nil
}

func loadConfigs(path string) ([]santa.Config, error) {
	var configs []santa.Config
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) != ".toml" {
			return nil
		}
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			var conf santa.Config
			err := toml.Unmarshal(file, &conf)
			if err != nil {
				return errors.Wrapf(err, "failed to decode %v, skipping \n", info.Name())
			}
			name := info.Name()
			conf.MachineID = strings.TrimSuffix(name, filepath.Ext(name))
			configs = append(configs, conf)
			return nil
		}
		return nil
	})
	return configs, errors.Wrapf(err, "loading configs from path")
}
