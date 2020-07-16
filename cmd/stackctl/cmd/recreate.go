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
	"bytes"
	"html/template"

	"github.com/compose-spec/compose-go/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/stackctl/pkg/actions"
	"github.com/trusch/stackctl/pkg/compose"
)

// recreateCmd represents the recreate command
var recreateCmd = &cobra.Command{
	Use:   "recreate",
	Short: "recreate your pod or just components",
	Long:  `recreate your pod or just components.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetStringSlice("compose-file")
		project, err := compose.Load(file)
		if err != nil {
			logrus.Fatal(err)
		}
		if len(args) == 1 {
			image, _ := cmd.Flags().GetString("image")
			pr, _ := cmd.Flags().GetString("with-pr")
			for idx, svc := range project.Services {
				if svc.Name == args[0] {
					if image != "" {
						project.Services[idx].Image = image
					}
					if pr != "" {
						project.Services[idx].Image = getPullRequestImage(project, svc, pr)
					}
				}
			}
		}
		err = actions.Recreate(cmd.Context(), project, args)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func getPullRequestImage(project *types.Project, svc types.ServiceConfig, pr string) string {
	tmplInterface, ok := svc.Extensions["x-pr-template"]
	if !ok {
		tmplInterface, ok = project.Extensions["x-pr-template"]
		if !ok {
			logrus.Fatal("no x-pr-template entry found")
		}
	}
	tmplStr, ok := tmplInterface.(string)
	if !ok {
		logrus.Fatal("x-pr-template is not a string")
	}
	tmpl, err := template.New("pr-template").Parse(tmplStr)
	if err != nil {
		logrus.Fatal(err)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct{ Service, PR string }{svc.Name, pr})
	if err != nil {
		logrus.Fatal(err)
	}
	return buf.String()
}

func init() {
	rootCmd.AddCommand(recreateCmd)
	recreateCmd.Flags().String("image", "", "override image")
	recreateCmd.Flags().String("with-pr", "", "override image by using the PR template")
}
