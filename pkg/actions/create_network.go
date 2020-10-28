package actions

import (
	"context"
	"errors"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func CreateNetworks(ctx context.Context, project *types.Project) error {
	for _, vol := range project.NetworkNames() {
		err := CreateNetwork(ctx, project, vol)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateNetwork(ctx context.Context, project *types.Project, network string) error {
	if err := checkIfNetworkExists(ctx, project.Name+"-"+network); err == nil {
		return nil
	}
	logrus.Infof("creating network %s", network)
	networkConfig, ok := project.Networks[network]
	if !ok {
		return errors.New("no such network")
	}
	args := []string{"network", "create"}
	if networkConfig.Driver != "" {
		args = append(args, "--driver", networkConfig.Driver)
	}
	if len(networkConfig.DriverOpts) > 0 {
		builder := &strings.Builder{}
		idx := 0
		for k, v := range networkConfig.DriverOpts {
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
	if len(networkConfig.Labels) > 0 {
		builder := &strings.Builder{}
		idx := 0
		for k, v := range networkConfig.DriverOpts {
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
	args = append(args, project.Name+"-"+network)
	return run(ctx, "podman", args...)
}

func checkIfNetworkExists(ctx context.Context, name string) error {
	args := []string{"network", "inspect", name}
	return run(ctx, "podman", args...)
}
