package cb

import (
	"reflect"

	"github.com/anarcher/cue-bundler/pkg/spec/v1/deps"
)

func depEqual(d1, d2 deps.Dependency) bool {
	name := d1.Name() == d2.Name()
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}
