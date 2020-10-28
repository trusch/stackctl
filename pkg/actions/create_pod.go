package actions

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
)

func CreatePod(ctx context.Context, project *types.Project) error {
	logrus.Info("creating pod")
	err := CreateVolumes(ctx, project)
	if err != nil {
		return err
	}

	err = CreateNetworks(ctx, project)
	if err != nil {
		return err
	}

	args := []string{"pod", "create", "--name", project.Name, "--share", "cgroup,ipc,uts"}
	/*
		  // add port mappings
		  services, err := project.GetServices(nil)
		  if err != nil {
		  	return err
		  }

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

			if portsInterface, ok := project.Extensions["x-ports"]; ok {
				if portsObj, ok := portsInterface.(map[string]interface{}); ok {
					for k, v := range portsObj {
						if target, ok := v.(string); ok {
							args = append(args, "-p", k+":"+target)
						}
					}
				}
			}
	*/

	return run(ctx, "podman", args...)
}
