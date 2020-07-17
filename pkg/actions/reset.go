package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
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
		logrus.Infof("removing volume %s", vol)
		args = append(args, vol)
	}
	return run(ctx, "podman", args...)
}
