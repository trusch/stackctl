package actions

import (
	"context"
	"strings"

	"github.com/compose-spec/compose-go/types"
)

func Up(ctx context.Context, project *types.Project) error {
	err := CreateVolumes(ctx, project)
	if err != nil {
		return err
	}

	status, err := GetStatus(ctx, project)
	if err != nil {
		return err
	}
	if len(status) == 0 {
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
	} else {
		for _, st := range status {
			if !strings.HasPrefix(st.Status, "Up") {
				err = StartService(ctx, project, st.Name)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
