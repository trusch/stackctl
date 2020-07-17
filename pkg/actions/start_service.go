package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func StartService(ctx context.Context, project *types.Project, svc string) error {
	logrus.Infof("starting service %s", svc)
	args := []string{"container", "start", containerName(project, svc)}
	return run(ctx, "podman", args...)
}
