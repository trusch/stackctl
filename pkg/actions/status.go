package actions

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/olekukonko/tablewriter"
)

type ProjectStatus []Status

type Status struct {
	ID     string
	Name   string
	Image  string
	Status string
}

func GetStatus(ctx context.Context, project *types.Project) (results ProjectStatus, err error) {
	args := []string{"ps", "-a", "--filter", "label=pod=" + project.Name, "--sort", "names", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}"}
	cmd := exec.CommandContext(ctx, "podman", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			results = append(results, Status{
				ID:     parts[0],
				Name:   parts[1],
				Image:  parts[2],
				Status: parts[3],
			})
		}
	}
	return results, nil
}

func (p ProjectStatus) Print() error {
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetBorder(false)
	writer.SetHeader([]string{"ID", "NAME", "IMAGE", "STATUS"})
	for _, status := range p {
		writer.Append([]string{status.ID, status.Name, status.Image, status.Status})
	}
	writer.Render()
	return nil
}
