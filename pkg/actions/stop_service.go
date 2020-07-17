package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func StopService(ctx context.Context, project *types.Project, service string) error {
	args := []string{"container", "stop", "-t", "1", containerName(project, service)}
	return run(ctx, "podman", args...)
}
