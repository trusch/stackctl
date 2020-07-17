package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func StopService(ctx context.Context, project *types.Project, service string) error {
	logrus.Infof("stopping service %s", service)
	args := []string{"container", "stop", "-t", "1", containerName(project, service)}
	return run(ctx, "podman", args...)
}
