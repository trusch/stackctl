package podman

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/trusch/stackctl/pkg/config"
)

func CreatePod(cfg config.Config) error {
	args := []string{"pod", "create", "--name", cfg.Name}
	for host, container := range cfg.Ports {
		args = append(args,
			"-p", fmt.Sprintf("%s:%s", host, container),
		)
	}
	err := run(args)
	if err != nil {
		return err
	}
	return createVolumesIfMissing(cfg)
}

func DeletePod(cfg config.Config) error {
	args := []string{"pod", "rm", "-f", cfg.Name}
	return run(args)
}

func StartPod(cfg config.Config) error {
	args := []string{"pod", "start", cfg.Name}
	return run(args)
}

func StopPod(cfg config.Config) error {
	args := []string{"pod", "stop", "-t", "1", cfg.Name}
	return run(args)
}

func DeleteVolumes(cfg config.Config) error {
	for name, src := range cfg.Volumes {
		if src == "" {
			args := []string{"volume", "rm", volumeName(cfg, name)}
			err := run(args)
			if err != nil {
				logrus.Warn(err)
			}
		}
	}
	return nil
}

func createVolumesIfMissing(cfg config.Config) error {
	for name, target := range cfg.Volumes {
		if target == "" {
			args := []string{"volume", "create", volumeName(cfg, name)}
			_ = runSilent(args)
		}
	}
	return nil
}

func volumeName(cfg config.Config, name string) string {
	return fmt.Sprintf("vol-%s-%s", cfg.Name, name)
}
