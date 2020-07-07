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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		cfg, err := config.Load(file)
		if err != nil {
			logrus.Fatal(err)
		}
		if len(args) == 0 {
			logrus.Infof("restarting pod")
			err = podman.StopPod(cfg)
			if err != nil {
				logrus.Warn(err)
			}
			err = podman.StartPod(cfg)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			for _, component := range cfg.Components {
				if args[0] == component.Name {
					logrus.Infof("restarting container for %s", component.Name)
					err = podman.StopContainer(cfg, component)
					if err != nil {
						logrus.Warn(err)
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
	rootCmd.AddCommand(restartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
