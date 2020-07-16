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
	"github.com/trusch/stackctl/pkg/actions"
	"github.com/trusch/stackctl/pkg/compose"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"down", "delete"},
	Short:   "remove the whole stack or just components",
	Long:    `remove the whole stack or just components.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetStringSlice("compose-file")
		project, err := compose.Load(file)
		if len(args) == 0 {
			logrus.Infof("removing pod")
			err = actions.RemovePod(cmd.Context(), project)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			for _, svc := range project.ServiceNames() {
				if args[0] == svc {
					logrus.Infof("removing service %s", svc)
					err = actions.RemoveService(cmd.Context(), project, svc)
					if err != nil {
						logrus.Fatal(err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
