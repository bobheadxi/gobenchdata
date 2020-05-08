package web

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

//go:generate npm run build
//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

func defaultConfigPath(dir string) string { return path.Join(dir, "gobenchdata-web.yml") }

// GenerateApp dumps the web app template into the provided directory
func GenerateApp(dir string, it TemplateIndexHTML) error {
	// clear directory of build stuff
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot access target directory: %w", err)
	}
	for _, f := range dirs {
		fullName := path.Join(dir, f.Name())
		if f.IsDir() {
			// clear build folders
			if f.Name() == "css" || f.Name() == "js" {
				os.RemoveAll(fullName)
			}
		} else if f.Name() == "index.html" {
			os.Remove(fullName)
		}
	}

	if err := populateFileIndexHTML(it); err != nil {
		return err
	}

	// generate app
	appFiles, _ := WalkDirs(".", false)
	for _, f := range appFiles {
		b, _ := ReadFile(f)
		target := path.Join(dir, f)
		os.MkdirAll(path.Dir(target), os.ModePerm)
		if err := ioutil.WriteFile(target, b, os.ModePerm); err != nil {
			return fmt.Errorf("failed to generate '%s': %w", f, err)
		}
	}

	return nil
}

// GenerateConfig generates configuration for the web app
func GenerateConfig(dir string, defaultConfig Config, override bool) error {
	appConfigPath := defaultConfigPath(dir)

	// check for existing
	if !override {
		if _, err := os.Stat(appConfigPath); err == nil {
			return fmt.Errorf("found existing configuration: %w", os.ErrExist)
		}
	} else {
		os.Remove(appConfigPath)
	}

	// generate configuration
	b, _ := yaml.Marshal(&defaultConfig)
	if err := ioutil.WriteFile(appConfigPath, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	return nil
}
