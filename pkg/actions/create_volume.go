package actions

import (
	"context"
	"errors"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func CreateVolumes(ctx context.Context, project *types.Project) error {
	for _, vol := range project.VolumeNames() {
		err := CreateVolume(ctx, project, vol)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateVolume(ctx context.Context, project *types.Project, volume string) error {
	if err := checkIfVolumeExists(ctx, volume); err == nil {
		return nil
	}
	logrus.Infof("creating volume %s", volume)
	volumeConfig, ok := project.Volumes[volume]
	if !ok {
		return errors.New("no such volume")
	}
	args := []string{"volume", "create"}
	if volumeConfig.Driver != "" {
		args = append(args, "--driver", volumeConfig.Driver)
	}
	if len(volumeConfig.DriverOpts) > 0 {
		builder := &strings.Builder{}
		idx := 0
		for k, v := range volumeConfig.DriverOpts {
			if idx > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(k)
			if v != "" {
				builder.WriteString("=")
				builder.WriteString(v)
			}
			idx++
		}
		args = append(args, "--opt", builder.String())
	}
	if len(volumeConfig.Labels) > 0 {
		builder := &strings.Builder{}
		idx := 0
		for k, v := range volumeConfig.DriverOpts {
			if idx > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(k)
			if v != "" {
				builder.WriteString("=")
				builder.WriteString(v)
			}
			idx++
		}
		args = append(args, "--label", builder.String())
	}
	return run(ctx, "podman", args...)
}

func checkIfVolumeExists(ctx context.Context, name string) error {
	args := []string{"volume", "inspect", name}
	return run(ctx, "podman", args...)
}
