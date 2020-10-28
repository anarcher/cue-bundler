package cb

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/anarcher/cue-bundler/pkg/spec"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
	"github.com/anarcher/cue-bundler/pkg/spec/v1/deps"
)

func Ensure(cfg v1.Config, modDir string, oldLocks map[string]deps.Dependency) (map[string]deps.Dependency, error) {
	pkgDir := filepath.Join(modDir, "pkg")

	locks, err := ensure(cfg.Dependencies, pkgDir, oldLocks)
	if err != nil {
		return nil, err
	}

	return locks, nil
}

func ensure(direct map[string]deps.Dependency, vendorDir string, locks map[string]deps.Dependency) (map[string]deps.Dependency, error) {
	deps := make(map[string]deps.Dependency)

	for _, d := range direct {
		l, present := locks[d.Name()]

		if present {
			d.Version = locks[d.Name()].Version

			if check(l, vendorDir) {
				deps[d.Name()] = l
				continue
			}
		}
		expectedSum := locks[d.Name()].Sum

		dir := filepath.Join(vendorDir, d.Name())
		os.RemoveAll(dir)

		locked, err := download(d, vendorDir)
		if err != nil {
			return nil, fmt.Errorf("downloading %w", err)
		}

		if expectedSum != "" && locked.Sum != expectedSum {
			return nil, fmt.Errorf("checksum mismatch for %s. Expected %s but got %s", d.Name(), expectedSum, locked.Sum)
		}

		deps[d.Name()] = *locked
		locks[d.Name()] = *locked
	}

	for _, d := range deps {
		if d.Single {
			continue
		}
		f, err := spec.Load(filepath.Join(vendorDir, d.Name(), spec.ConfigFile))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		nested, err := ensure(f.Dependencies, vendorDir, locks)
		if err != nil {
			return nil, err
		}

		for _, d := range nested {
			if _, ok := deps[d.Name()]; !ok {
				deps[d.Name()] = d
			}
		}
	}

	return deps, nil
}

// download retrieves a package from a remote upstream. The checksum of the
// files is generated afterwards.
func download(d deps.Dependency, vendorDir string) (*deps.Dependency, error) {
	var p Interface
	switch {
	case d.Source.GitSource != nil:
		p = NewGitPackage(d.Source.GitSource)
	case d.Source.LocalSource != nil:
		p = NewLocalPackage(d.Source.LocalSource)
	}

	if p == nil {
		return nil, errors.New("either git or local source is required")
	}

	version, err := p.Install(context.TODO(), d.Name(), vendorDir, d.Version)
	if err != nil {
		return nil, err
	}

	var sum string
	if d.Source.LocalSource == nil {
		sum = hashDir(filepath.Join(vendorDir, d.Name()))
	}

	d.Version = version
	d.Sum = sum
	return &d, nil
}

func check(d deps.Dependency, vendorDir string) bool {
	// assume a local dependency is intact as long as it exists
	if d.Source.LocalSource != nil {
		x, err := spec.Exists(filepath.Join(vendorDir, d.Name()))
		if err != nil {
			return false
		}
		return x
	}

	if d.Sum == "" {
		// no sum available, need to download
		return false
	}

	dir := filepath.Join(vendorDir, d.Name())
	sum := hashDir(dir)
	return d.Sum == sum

}

// hashDir computes the checksum of a directory by concatenating all files and
// hashing this data using sha256. This can be memory heavy with lots of data,
// but jsonnet files should be fairly small
func hashDir(dir string) string {
	hasher := sha256.New()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(hasher, f); err != nil {
			return err
		}

		return nil
	})

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
