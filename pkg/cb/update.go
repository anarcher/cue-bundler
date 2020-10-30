package cb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anarcher/cue-bundler/pkg/spec"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
	"github.com/anarcher/cue-bundler/pkg/spec/v1/deps"
)

func Update(modDir string, uris []string) error {
	l, err := load(modDir)
	if err != nil {
		return err
	}

	oldloaded := l.Clone()

	locks := l.lock.Dependencies

	//TODO(anarcher):
	if err := os.MkdirAll(filepath.Join(modDir, "pkg", ".tmp"), os.ModePerm); err != nil {
		return err
	}

	for _, u := range uris {
		d := deps.Parse(modDir, u)
		if d == nil {
			//TODO(anarcher):
			return fmt.Errorf("Unable to parse package URI %s", u)
		}

		delete(locks, d.Name())
	}

	if len(locks) == 0 {
		locks = make(map[string]deps.Dependency)
	}

	locked, err := Ensure(l.cfg, modDir, locks)
	if err != nil {
		return err
	}

	// write
	if err := WriteChangedConfigFile(oldloaded.lock, v1.Config{Version: v1.Version, Dependencies: locked}, filepath.Join(modDir, spec.ConfigLockFile)); err != nil {
		return err
	}

	return nil
}
