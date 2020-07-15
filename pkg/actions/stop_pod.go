package actions

import (
	"context"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/types"
)

func StopPod(ctx context.Context, project *types.Project) error {
	args := []string{"pod", "stop", "-t", "1", project.Name}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
