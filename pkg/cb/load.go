package cb

import (
	"path/filepath"

	"github.com/anarcher/cue-bundler/pkg/spec"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
)

type loaded struct {
	modDir string
	cfg    v1.Config
	lock   v1.Config
}

func load(modDir string) (*loaded, error) {
	cfg, err := spec.Load(filepath.Join(modDir, spec.ConfigFile))
	if err != nil {
		return nil, err
	}
	lock, err := spec.Load(filepath.Join(modDir, spec.ConfigLockFile))
	if err != nil {
		return nil, err
	}

	c := &loaded{
		modDir: modDir,
		cfg:    cfg,
		lock:   lock,
	}

	return c, nil
}

func (l loaded) Clone() *loaded {
	n := &loaded{
		modDir: l.modDir,
		cfg:    l.cfg,
		lock:   l.lock,
	}
	return n
}
