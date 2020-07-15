package actions

import (
	"context"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/types"
)

func Pull(ctx context.Context, project *types.Project, svc string) error {
	svcConfig, err := project.GetService(svc)
	if err != nil {
		return err
	}
	args := []string{"pull", svcConfig.Image}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
