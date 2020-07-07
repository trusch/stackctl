/*
Copyright © 2020 Tino Rusch <tino.rusch@contiamo.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/stackctl/pkg/config"
	"github.com/trusch/stackctl/pkg/podman"
)

// recreateCmd represents the recreate command
var recreateCmd = &cobra.Command{
	Use:   "recreate",
	Short: "recreate your pod or just components",
	Long:  `recreate your pod or just components.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		cfg, err := config.Load(file)
		if err != nil {
			logrus.Fatal(err)
		}
		if len(args) == 0 {
			logrus.Infof("recreating pod")
			_ = podman.StopPod(cfg)
			_ = podman.DeletePod(cfg)
			err = podman.CreatePod(cfg)
			if err != nil {
				logrus.Fatal(err)
			}
			for _, component := range cfg.Components {
				err = podman.CreateContainer(cfg, component)
				if err != nil {
					logrus.Fatal(err)
				}
			}
			err = podman.StartPod(cfg)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			image, _ := cmd.Flags().GetString("image")
			for _, component := range cfg.Components {
				if args[0] == component.Name {
					logrus.Infof("recreating container for %s", component.Name)
					_ = podman.StopContainer(cfg, component)
					_ = podman.DeleteContainer(cfg, component)
					if image != "" {
						component.Image = image
					}
					err = podman.CreateContainer(cfg, component)
					if err != nil {
						logrus.Fatal(err)
					}
					err = podman.StartContainer(cfg, component)
					if err != nil {
						logrus.Fatal(err)
					}
					break
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(recreateCmd)
	recreateCmd.Flags().String("image", "", "override image")
}
