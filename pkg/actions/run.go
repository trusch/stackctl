package actions

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func run(ctx context.Context, bin string, args ...string) error {
	logrus.Debugf("run> %s %s", bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = os.Stdin
	if logrus.GetLevel() == logrus.DebugLevel {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
