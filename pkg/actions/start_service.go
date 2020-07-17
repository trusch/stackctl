package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func StartService(ctx context.Context, project *types.Project, svc string) error {
	args := []string{"container", "start", containerName(project, svc)}
	return run(ctx, "podman", args...)
}
