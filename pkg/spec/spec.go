package spec

import (
	"errors"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	v1 "github.com/anarcher/cue-bundler/pkg/spec/v1"
)

const ConfigFile = "cb.cue"
const ConfigLockFile = "cb.lock.cue"

var ErrConfigVersion = errors.New("version unknown, update cb")

func Load(filepath string) (v1.Config, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return v1.New(), err
	}
	return Unmarshal(bs)
}

// Exists returns whether the file at the given path exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func Unmarshal(bs []byte) (v1.Config, error) {
	m := v1.New()

	if len(bs) == 0 {
		return m, nil
	}

	var r cue.Runtime
	inst, err := r.Compile("config", bs)
	if err != nil {
		return m, err
	}

	version, err := inst.Lookup("version").String()
	if err != nil {
		return m, err
	}

	switch version {
	case v1.Version:
		if err := inst.Value().Decode(&m); err != nil {
			return m, err
		}
	default:
		return m, ErrConfigVersion
	}

	return m, nil
}
