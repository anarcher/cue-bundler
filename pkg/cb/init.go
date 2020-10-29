package cb

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/anarcher/cue-bundler/pkg/cueutil"
	"github.com/anarcher/cue-bundler/pkg/spec"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
)

func Init(modDir string) error {
	if err := InitConfigFile(modDir, spec.ConfigFile); err != nil {
		return err
	}
	if err := InitConfigFile(modDir, spec.ConfigLockFile); err != nil {
		return err
	}

	return nil
}

func InitConfigFile(modDir, filename string) error {
	path := filepath.Join(modDir, filename)
	exists, err := spec.Exists(path)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s already exists", path)
	}

	cfg := v1.New()
	bs, err := cueutil.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bs, 0644)
}
