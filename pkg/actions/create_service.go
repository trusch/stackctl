package actions

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/compose-spec/compose-go/types"
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

	// add volumes
	for _, vol := range append(svc.Volumes) {
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

	// add pod label
	args = append(args, "--label", "pod="+project.Name)

	// set image
	args = append(args, svc.Image)

	// set command
	if len(svc.Command) > 0 {
		args = append(args, svc.Command...)
	}

	// run it
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func containerName(project *types.Project, service string) string {
	return project.Name + "-" + service
}
