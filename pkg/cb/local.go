package cb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anarcher/cue-bundler/pkg/spec/v1/deps"
	"github.com/fatih/color"
)

type LocalPackage struct {
	Source *deps.Local
}

func NewLocalPackage(source *deps.Local) Interface {
	return &LocalPackage{
		Source: source,
	}
}

func (p *LocalPackage) Install(ctx context.Context, name, dir, version string) (lockVersion string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	oldname := filepath.Join(wd, p.Source.Directory)
	newname := filepath.Join(dir, name)
	linkname, err := filepath.Rel(dir, oldname)
	if err != nil {
		linkname = oldname
	}

	err = os.RemoveAll(newname)
	if err != nil {
		return "", fmt.Errorf("failed to clean previous destination path: %w", err)
	}

	_, err = os.Stat(oldname)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("symlink destination path does not exist: %w", err)
	}

	err = os.Symlink(linkname, newname)
	if err != nil {
		return "", fmt.Errorf("failed to create symlink for local dependency: %w", err)
	}

	color.Magenta("LOCAL %s -> %s", name, oldname)
	return "", nil
}
