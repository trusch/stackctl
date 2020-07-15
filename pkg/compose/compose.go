package compose

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

func Load(file string) (project *types.Project, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	contents, err := loader.ParseYAML(bs)
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	project, err = loader.Load(types.ConfigDetails{
		Version:    "3.9",
		WorkingDir: wd,
		ConfigFiles: []types.ConfigFile{
			{
				Filename: file,
				Config:   contents,
			},
		},
		Environment: map[string]string{},
	})
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	project.Name = filepath.Base(filepath.Dir(path))

	return project, nil
}
