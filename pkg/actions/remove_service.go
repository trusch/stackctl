package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func RemoveService(ctx context.Context, project *types.Project, service string) error {
	args := []string{"container", "rm", "-f", containerName(project, service)}
	return run(ctx, "podman", args...)
}
