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
		st, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		switch {
		case st.IsDir():
			configs, err := parseDirectory(f)
			if err != nil {
				return nil, err
			}
			configDetails.ConfigFiles = append(configDetails.ConfigFiles, configs...)
		case !st.IsDir():
			config, err := parseFile(f)
			if err != nil {
				return nil, err
			}
			configDetails.ConfigFiles = append(configDetails.ConfigFiles, config)
		}
	}
	project, err = loader.Load(configDetails)
	if err != nil {
		return nil, err
	}
	project.Name = filepath.Base(wd)

	return project, nil
}

func parseFile(file string) (types.ConfigFile, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return types.ConfigFile{}, err
	}
	contents, err := loader.ParseYAML(bs)
	if err != nil {
		return types.ConfigFile{}, err
	}
	return types.ConfigFile{
		Filename: file,
		Config:   contents,
	}, nil
}

func parseDirectory(dir string) (result []types.ConfigFile, err error) {
	return result, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			config, err := parseFile(path)
			if err != nil {
				return err
			}
			result = append(result, config)
		}
		return nil
	})
}
