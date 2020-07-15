package actions

import (
	"context"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/types"
)

func StartService(ctx context.Context, project *types.Project, svc string) error {
	args := []string{"container", "start", containerName(project, svc)}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
