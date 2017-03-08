package santaconfig

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/groob/moroz/santa"
)

func NewFileRepo(path string) *FileRepo {
	repo := FileRepo{
		configIndex: make(map[string]*santa.Config),
		configPath:  path,
	}
	return &repo
}

type FileRepo struct {
	mtx         sync.RWMutex
	configIndex map[string]*santa.Config
	configPath  string
}

func (f *FileRepo) updateIndex(configs *santa.ConfigCollection) {
	f.configIndex = make(map[string]*santa.Config, len(*configs))
	for _, conf := range *configs {
		f.configIndex[conf.MachineID] = conf
	}
}

func (f *FileRepo) AllConfigs() (*santa.ConfigCollection, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	var configs santa.ConfigCollection
	err := loadConfigs(f.configPath, &configs)
	if err != nil {
		return nil, err
	}
	f.updateIndex(&configs)
	return &configs, nil
}

func (f *FileRepo) Config(machineID string) (*santa.Config, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	var configs santa.ConfigCollection
	err := loadConfigs(f.configPath, &configs)
	if err != nil {
		return nil, err
	}
	f.updateIndex(&configs)
	conf, ok := f.configIndex[machineID]
	if !ok {
		return nil, errors.New("configuration not found")
	}
	return conf, nil
}

func loadConfigs(path string, configs *santa.ConfigCollection) error {
	return filepath.Walk(path, walkConfigs(configs))
}

func walkConfigs(configs *santa.ConfigCollection) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
				log.Printf("failed to decode %v, skipping \n", info.Name())
				return nil
			}
			name := info.Name()
			conf.MachineID = strings.TrimSuffix(name, filepath.Ext(name))
			*configs = append(*configs, &conf)
			return nil
		}
		return nil
	}
}
