package actions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func CreateServices(ctx context.Context, project *types.Project) error {
	for _, name := range project.ServiceNames() {
		err := CreateService(ctx, project, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateService(ctx context.Context, project *types.Project, name string) error {
	logrus.Infof("creating service %s", name)
	svc, err := project.GetService(name)
	if err != nil {
		return err
	}
	args := []string{"container", "create", "--pod", project.Name, "--name", containerName(project, name), "--add-host", "localhost:127.0.0.1"}

	// add env
	for k, v := range svc.Environment {
		if v == nil {
			args = append(args, "-e", k+"=")
		} else {
			args = append(args, "-e", k+"="+*v)
		}
	}

	// add networks
	for net := range svc.Networks {
		args = append(args, "--net", project.Name+"-"+net)
	}

	args = append(args, "--network-alias", name)

	// add volumes
	for _, vol := range svc.Volumes {
		args = append(args, "-v", vol.Source+":"+vol.Target)
	}

	// add configs
	for _, vol := range append(svc.Configs) {
		source := project.Configs[vol.Source].File
		target := vol.Target
		if target == "" {
			target = "/" + vol.Source
		}
		args = append(args, "-v", source+":"+target)
	}

	// add secrets
	for _, vol := range append(svc.Secrets) {
		source := project.Secrets[vol.Source].File
		target := vol.Target
		if target == "" {
			target = "/run/secrets/" + vol.Source
		}
		args = append(args, "-v", source+":"+target)
	}

	// set entrypoint
	if len(svc.Entrypoint) > 0 {
		args = append(args, "--entrypoint", strings.Join(svc.Entrypoint, " "))
	}

	// add extra hosts
	if len(svc.ExtraHosts) > 0 {
		args = append(args, "--add-host", strings.Join(svc.ExtraHosts, ","))
	}

	// set running user
	if svc.User != "" {
		args = append(args, "--user", svc.User)
	}

	// set working directory
	if svc.WorkingDir != "" {
		args = append(args, "-w", svc.WorkingDir)
	}

	// add labels
	if len(svc.Labels) > 0 {
		for k, v := range svc.Labels {
			args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
		}
	}

	// set port forwardings
	for _, port := range svc.Ports {
		builder := &strings.Builder{}
		if port.HostIP != "" {
			builder.WriteString(port.HostIP)
			builder.WriteString(":")
		}
		builder.WriteString(strconv.Itoa(int(port.Published)))
		builder.WriteString(":")
		builder.WriteString(strconv.Itoa(int(port.Target)))
		if port.Protocol != "" {
			builder.WriteString("/")
			builder.WriteString(port.Protocol)
		}
		args = append(args, "-p", builder.String())
	}

	// add pod label
	args = append(args, "--label", "pod="+project.Name)

	for name, limit := range svc.Ulimits {
		if limit.Single > 0 {
			args = append(args, "--ulimit", fmt.Sprintf("%v=%v", name, limit.Single))
		} else {
			args = append(args, "--ulimit", fmt.Sprintf("%v=%v:%v", name, limit.Soft, limit.Hard))
		}
	}

	// set image
	args = append(args, svc.Image)

	// set command
	if len(svc.Command) > 0 {
		args = append(args, svc.Command...)
	}

	// run it
	return run(ctx, "podman", args...)
}

func containerName(project *types.Project, service string) string {
	return project.Name + "-" + service
}
