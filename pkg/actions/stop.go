package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func Stop(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) == 0 {
		err := StopPod(ctx, project)
		if err != nil {
			return err
		}
		return nil
	}
	for _, svc := range svcs {
		err := StopService(ctx, project, svc)
		if err != nil {
			logrus.Error(err)
		}
	}
	return nil
}
