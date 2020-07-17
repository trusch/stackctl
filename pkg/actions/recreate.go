package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func Recreate(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) == 0 {
		err := StopPod(ctx, project)
		if err != nil {
			logrus.Warn(err)
		}
		err = RemovePod(ctx, project)
		if err != nil {
			logrus.Warn(err)
		}
		err = CreatePod(ctx, project)
		if err != nil {
			return err
		}
		for _, svc := range project.ServiceNames() {
			err = CreateService(ctx, project, svc)
			if err != nil {
				return err
			}
		}
		err = StartPod(ctx, project)
		if err != nil {
			return err
		}
		return nil
	}
	for _, svc := range svcs {
		err := StopService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
		err = RemoveService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
		err = CreateService(ctx, project, svc)
		if err != nil {
			return err
		}
		err = StartService(ctx, project, svc)
		if err != nil {
			return err
		}
	}
	return nil
}
