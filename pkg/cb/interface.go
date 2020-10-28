package cb

import "context"

type Interface interface {
	Install(ctx context.Context, name, dir, version string) (lockVersion string, err error)
}
