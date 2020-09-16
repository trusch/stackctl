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

// unpackSnapshotCmd represents the unpackSnapshot command
var unpackSnapshotCmd = &cobra.Command{
	Use:   "unpack",
	Short: "unpack a snapshot",
	Long:  `unpack a snapshot.`,
	Run: func(cmd *cobra.Command, args []string) {
		pw, _ := cmd.Flags().GetString("password")
		output, _ := cmd.Flags().GetString("output")
		input, _ := cmd.Flags().GetString("input")
		err := snapshot.Unpack(cmd.Context(), snapshot.UnpackOptions{
			OutputRoot:   output,
			Password:     pw,
			SnapshotFile: input,
		})
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	snapshotCmd.AddCommand(unpackSnapshotCmd)
	unpackSnapshotCmd.Flags().String("password", "", "password for decryption")
	unpackSnapshotCmd.Flags().String("input", "dev.snapshot", "input file")
	unpackSnapshotCmd.Flags().String("output", ".", "output root")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unpackSnapshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unpackSnapshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
