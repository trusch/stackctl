package podman

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/trusch/stackctl/pkg/config"
)

func Status(cfg config.Config) error {
	args := []string{"ps", "-a", "--filter", "label=pod=" + cfg.Name, "--format", "{{.ID}} {{.Names}} {{.Image}} {{.Status}}"}
	output, err := runAndGetOutput(args)
	if err != nil {
		return err
	}
	lines := strings.Split(string(output), "\n")
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetBorder(false)
	writer.SetHeader([]string{"ID", "NAME", "IMAGE", "STATUS"})
	for _, line := range lines {
		if line != "" {
			parts := strings.SplitN(line, " ", 4)
			parts[1] = parts[1][len(cfg.Name)+1:]
			writer.Append(parts)
		}
	}
	writer.Render()
	return nil
}
