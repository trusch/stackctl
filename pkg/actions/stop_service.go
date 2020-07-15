package actions

import (
	"context"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/types"
)

func StopService(ctx context.Context, project *types.Project, service string) error {
	args := []string{"container", "stop", "-t", "1", containerName(project, service)}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
