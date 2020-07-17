package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func StopPod(ctx context.Context, project *types.Project) error {
	args := []string{"pod", "stop", "-t", "1", project.Name}
	return run(ctx, "podman", args...)
}
