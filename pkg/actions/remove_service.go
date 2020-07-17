package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func RemoveService(ctx context.Context, project *types.Project, service string) error {
	logrus.Infof("removing service %s", service)
	args := []string{"container", "rm", "-f", containerName(project, service)}
	return run(ctx, "podman", args...)
}
