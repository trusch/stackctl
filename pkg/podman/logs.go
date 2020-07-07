package podman

import (
	"strconv"
	"time"

	"github.com/trusch/stackctl/pkg/config"
)

func Logs(cfg config.Config, components []string, since time.Time, tail int, follow bool) error {
	args := []string{"logs"}
	if len(components) > 1 {
		args = append(args, "--names")
	}
	if !since.IsZero() {
		args = append(args, "--since", since.Format(time.RFC3339))
	}
	if tail > 0 {
		args = append(args, "--tail", strconv.Itoa(tail))
	}
	if follow {
		args = append(args, "--follow")
	}
	for _, component := range components {
		args = append(args, containerName(cfg, config.Component{
			Name: component,
		}))
	}
	return run(args)
}
