package actions

import (
	"context"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/types"
)

func Reset(ctx context.Context, project *types.Project) error {
	_ = StopPod(ctx, project)
	_ = RemovePod(ctx, project)
	if len(project.Volumes) > 0 {
		return RemoveVolumes(ctx, project)
	}
	return nil
}

func RemoveVolumes(ctx context.Context, project *types.Project) error {
	args := []string{"volume", "rm"}
	for _, vol := range project.VolumeNames() {
		args = append(args, vol)
	}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
