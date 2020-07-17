package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func StartPod(ctx context.Context, project *types.Project) error {
	logrus.Infof("starting pod")
	args := []string{"pod", "start", project.Name}
	return run(ctx, "podman", args...)
}
