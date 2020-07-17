package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func Remove(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) > 0 {
		for _, svc := range svcs {
			_ = StopService(ctx, project, svc)
			err := RemoveService(ctx, project, svc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	_ = StopPod(ctx, project)
	err := RemovePod(ctx, project)
	if err != nil {
		return err
	}
	return nil
}
