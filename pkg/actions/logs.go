package actions

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/compose-spec/compose-go/types"
)

func Logs(ctx context.Context, project *types.Project, services []string, since time.Time, tail int, follow bool) error {
	args := []string{"logs"}
	if len(services) > 1 {
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
	for _, service := range services {
		args = append(args, containerName(project, service))
	}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
