package compose

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

func Load(files []string) (project *types.Project, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configDetails := types.ConfigDetails{
		Version:     "3.9",
		WorkingDir:  wd,
		ConfigFiles: []types.ConfigFile{},
		Environment: map[string]string{},
	}

	for _, f := range files {
		bs, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		contents, err := loader.ParseYAML(bs)
		if err != nil {
			return nil, err
		}
		configDetails.ConfigFiles = append(configDetails.ConfigFiles, types.ConfigFile{
			Filename: f,
			Config:   contents,
		})
	}
	project, err = loader.Load(configDetails)
	if err != nil {
		return nil, err
	}
	project.Name = filepath.Base(wd)

	return project, nil
}
