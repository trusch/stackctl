package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

func Create(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) == 0 {
		err := CreatePod(ctx, project)
		if err != nil {
			return err
		}
		svcs = project.ServiceNames()
	}
	for _, svc := range svcs {
		err := CreateService(ctx, project, svc)
		if err != nil {
			return err
		}
	}
	return nil
}
