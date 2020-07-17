package actions

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/compose-spec/compose-go/types"
)

func Stats(ctx context.Context, project *types.Project, svcs []string, watch bool) error {
	args := []string{"container", "stats"}
	if !watch {
		args = append(args, "--no-stream", "--no-reset")
	}

	if len(svcs) == 0 {
		containerStati, err := GetStatus(ctx, project)
		if err != nil {
			return err
		}
		for _, status := range containerStati {
			if strings.HasPrefix(status.Status, "Up") {
				args = append(args, containerName(project, status.Name))
			}
		}
	} else {
		for _, svc := range svcs {
			args = append(args, containerName(project, svc))
		}
	}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
