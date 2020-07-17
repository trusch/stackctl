package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func StartPod(ctx context.Context, project *types.Project) error {
	args := []string{"pod", "start", project.Name}
	return run(ctx, "podman", args...)
}
