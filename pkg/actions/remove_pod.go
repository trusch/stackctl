package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func RemovePod(ctx context.Context, project *types.Project) error {
	args := []string{"pod", "rm", "-f", project.Name}
	return run(ctx, "podman", args...)
}
