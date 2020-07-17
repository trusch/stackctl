package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func Start(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) > 0 {
		for _, svc := range svcs {
			err := StartService(ctx, project, svc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := StartPod(ctx, project)
	if err != nil {
		return err
	}
	return nil
}
