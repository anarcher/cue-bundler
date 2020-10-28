package v1

import "github.com/anarcher/cue-bundler/pkg/spec/v1/deps"

const Version = "v1"

type Config struct {
	Version      string                     `json:"version"`
	Dependencies map[string]deps.Dependency `json:"dependencies"`
}

func New() Config {
	return Config{
		Version:      Version,
		Dependencies: make(map[string]deps.Dependency),
	}
}
