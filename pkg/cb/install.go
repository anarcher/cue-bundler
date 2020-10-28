package cb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"cuelang.org/go/cue"
	"github.com/anarcher/cue-bundler/pkg/spec"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
	"github.com/anarcher/cue-bundler/pkg/spec/v1/deps"
)

func Install(modDir string, uris []string) error {
	l, err := load(modDir)
	if err != nil {
		return err
	}

	oldloaded := l.Clone()

	//TODO(anarcher):
	if err := os.MkdirAll(filepath.Join(modDir, ".tmp"), os.ModePerm); err != nil {
		return err
	}

	for _, u := range uris {
		d := deps.Parse(modDir, u)
		if d == nil {
			//TODO(anarcher):
			return fmt.Errorf("Unable to parse package URI %s", u)
		}

		if !depEqual(l.cfg.Dependencies[d.Name()], *d) {
			l.cfg.Dependencies[d.Name()] = *d
		}

		delete(l.lock.Dependencies, d.Name())
	}

	locked, err := Ensure(l.cfg, modDir, l.lock.Dependencies)
	if err != nil {
		return err
	}

	// write Cfg and Lock
	if err := writeChangedConfigFile(oldloaded.cfg, l.cfg, filepath.Join(modDir, spec.ConfigFile)); err != nil {
		return err
	}
	fmt.Println("updating cb.cue")
	if err := writeChangedConfigFile(oldloaded.lock, v1.Config{Dependencies: locked}, filepath.Join(modDir, spec.ConfigLockFile)); err != nil {
		return err
	}
	fmt.Println("updating cb.lock.cue")

	return nil
}

func writeChangedConfigFile(oldcfg v1.Config, cfg v1.Config, path string) error {
	if reflect.DeepEqual(oldcfg, cfg) {
		return nil
	}

	var r cue.Runtime
	inst, err := r.Compile("cfg", "")
	if err != nil {
		return err
	}
	inst, err = inst.Fill(cfg)
	if err != nil {
		return err
	}

	bs, err := r.Marshal(inst)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bs, 0644)
}
