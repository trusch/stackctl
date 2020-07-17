package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func RemovePod(ctx context.Context, project *types.Project) error {
	logrus.Info("removing pod")
	args := []string{"pod", "rm", "-f", project.Name}
	return run(ctx, "podman", args...)
}
