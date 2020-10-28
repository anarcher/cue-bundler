package deps

import (
	"os"
	"path/filepath"
)

type Source struct {
	GitSource   *Git   `json:"git,omitempty"`
	LocalSource *Local `json:"local,omitempty"`
}

type Local struct {
	Directory string `json:"directory"`
}

func (s Source) Name() string {
	switch {
	case s.GitSource != nil:
		return s.GitSource.Name()
	case s.LocalSource != nil:
		return s.LocalSource.Name()
	default:
		return ""
	}
}

func (l Local) Name() string {
	p, err := filepath.Abs(l.Directory)
	if err != nil {
		panic("unable to create absolute path from local source directory: " + err.Error())
	}
	return filepath.Base(p)
}

func parseLocal(dir, p string) *Dependency {
	clean := filepath.Clean(p)
	abs := filepath.Join(dir, clean)

	info, err := os.Stat(abs)
	if err != nil {
		return nil
	}

	if !info.IsDir() {
		return nil
	}

	return &Dependency{
		Source: Source{
			LocalSource: &Local{
				Directory: clean,
			},
		},
		Version: "",
	}
}
