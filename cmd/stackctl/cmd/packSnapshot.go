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
	"github.com/trusch/stackctl/pkg/snapshot"
)

// packSnapshotCmd represents the packSnapshot command
var packSnapshotCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack a snapshot of your applications state",
	Long:  `pack a snapshot of your applications state.`,
	Run: func(cmd *cobra.Command, args []string) {
		pw, _ := cmd.Flags().GetString("password")
		output, _ := cmd.Flags().GetString("output")

		if len(args) > 0 {
			err := snapshot.Pack(cmd.Context(), snapshot.PackOptions{
				Output:   output,
				Password: pw,
				Sources:  args,
			})
			if err != nil {
				logrus.Fatal(err)
			}
			return
		}
		// logrus.Info("snapshotting project volumes")
		// file, _ := cmd.Flags().GetStringSlice("compose-file")
		// project, err := compose.Load(file)
		// if err != nil {
		// 	logrus.Fatal(err)
		// }
		// vols := make([]string, 0, len(project.Volumes)+len(project.Configs)+len(project.Secrets))
		// vols = append(vols, project.VolumeNames()...)
		// vols = append(vols, project.ConfigNames()...)
		// vols = append(vols, project.SecretNames()...)
		// err = snapshot.CreateFromVolumes(cmd.Context(), project, vols, output)
		// if err != nil {
		// 	logrus.Fatal(err)
		// }
	},
}

func init() {
	snapshotCmd.AddCommand(packSnapshotCmd)
	packSnapshotCmd.Flags().String("password", "", "password for encryption")
	packSnapshotCmd.Flags().String("output", "dev.snapshot", "output file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packSnapshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packSnapshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
