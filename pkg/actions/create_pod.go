package actions

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/types"
)

func CreatePod(ctx context.Context, project *types.Project) error {
	err := CreateVolumes(ctx, project)
	if err != nil {
		return err
	}

	args := []string{"pod", "create", "--name", project.Name}
	services, err := project.GetServices(nil)
	if err != nil {
		return err
	}

	// add port mappings
	for _, svc := range services {
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
	}

	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
