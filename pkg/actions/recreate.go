package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func Recreate(ctx context.Context, project *types.Project, svcs []string) error {
	if len(svcs) == 0 {
		logrus.Info("stop pod")
		err := StopPod(ctx, project)
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Info("remove pod")
		err = RemovePod(ctx, project)
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Info("create pod")
		err = CreatePod(ctx, project)
		if err != nil {
			return err
		}
		for _, svc := range project.ServiceNames() {
			logrus.Infof("create service %s", svc)
			err = CreateService(ctx, project, svc)
			if err != nil {
				return err
			}
		}
		logrus.Info("start pod")
		err = StartPod(ctx, project)
		if err != nil {
			return err
		}
		return nil
	}
	for _, svc := range svcs {
		logrus.Infof("stop service %s", svc)
		err := StopService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Infof("remove service %s", svc)
		err = RemoveService(ctx, project, svc)
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Infof("create service %s", svc)
		err = CreateService(ctx, project, svc)
		if err != nil {
			return err
		}
		logrus.Infof("start service %s", svc)
		err = StartService(ctx, project, svc)
		if err != nil {
			return err
		}
	}
	return nil
}
