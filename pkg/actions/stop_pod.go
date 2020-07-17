package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func StopPod(ctx context.Context, project *types.Project) error {
	logrus.Info("stopping pod")
	args := []string{"pod", "stop", "-t", "1", project.Name}
	return run(ctx, "podman", args...)
}
