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
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/stackctl/pkg/actions"
	"github.com/trusch/stackctl/pkg/compose"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [flags] COMPONENT [COMPONENT...]",
	Short: "view logs",
	Long:  `view logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			since  time.Time
			tail   int
			follow bool
			err    error
		)
		file, _ := cmd.Flags().GetStringSlice("compose-file")
		project, err := compose.Load(file)
		if err != nil {
			logrus.Fatal(err)
		}
		sinceStr, _ := cmd.Flags().GetString("since")
		if sinceStr != "" {
			since, err = time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				logrus.Fatal(err)
			}
		}
		tail, _ = cmd.Flags().GetInt("tail")
		follow, _ = cmd.Flags().GetBool("follow")
		err = actions.Logs(cmd.Context(), project, args, since, tail, follow)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().String("since", "", "show logs since timestamp")
	logsCmd.Flags().IntP("tail", "t", 0, "show latest N lines from the end of the log")
	logsCmd.Flags().BoolP("follow", "f", false, "follow new logs")
}
