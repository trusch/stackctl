package podman

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/trusch/stackctl/pkg/config"
)

func CreateContainer(cfg config.Config, component config.Component) error {
	args := []string{
		"create",
		"--pod", cfg.Name,
		"--label", "pod=" + cfg.Name,
		"--name", containerName(cfg, component),
	}

	if component.Entrypoint != "" {
		args = append(args,
			"--entrypoint", component.Entrypoint,
		)
	}

	for key, value := range component.Environment {
		args = append(args,
			"-e", fmt.Sprintf("%s=%s", key, value),
		)
	}

	for src, target := range component.Mounts {
		vol, ok := cfg.Volumes[src]
		if !ok {
			return fmt.Errorf("no such volume: %s", src)
		}
		if vol == "" {
			vol = volumeName(cfg, src)
		}
		args = append(args,
			"-v", fmt.Sprintf("%s:%s", vol, target),
		)
	}

	args = append(args, "--add-host", "localhost:127.0.0.1")
	args = append(args, "--add-host", component.Name+":127.0.0.1")
	for host, ip := range component.Hosts {
		args = append(args,
			"--add-host", fmt.Sprintf("%s:%s", host, ip),
		)
	}

	args = append(args, component.Image)
	args = append(args, component.Cmd...)

	return run(args)
}

func DeleteContainer(cfg config.Config, component config.Component) error {
	args := []string{"rm", "-f", containerName(cfg, component)}
	return run(args)
}

func PullImage(cfg config.Config, component config.Component) error {
	args := []string{"pull", component.Image}
	return run(args)
}

func StartContainer(cfg config.Config, component config.Component) error {
	args := []string{"start", containerName(cfg, component)}
	return run(args)
}

func StopContainer(cfg config.Config, component config.Component) error {
	args := []string{"stop", "-t", "1", containerName(cfg, component)}
	return run(args)
}

func containerName(cfg config.Config, component config.Component) string {
	return fmt.Sprintf("%s-%s", cfg.Name, component.Name)
}

func run(args []string) error {
	cmd := exec.Command("podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Debugf("run> podman %v", args)
	return cmd.Run()
}

func runAndGetOutput(args []string) ([]byte, error) {
	buf := &bytes.Buffer{}
	cmd := exec.Command("podman", args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	logrus.Debugf("run> podman %v", args)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func runSilent(args []string) error {
	cmd := exec.Command("podman", args...)
	logrus.Debugf("run silent> podman %v", args)
	return cmd.Run()
}
