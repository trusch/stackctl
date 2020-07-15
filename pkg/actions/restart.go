package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func Restart(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) == 0 {
		err := StopPod(ctx, project)
		if err != nil {
			logrus.Warn(err)
		}
		err = StartPod(ctx, project)
		return err
	}
	for _, svc := range svcs {
		err := StopService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
		err = StartService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
	}
	return nil
}
