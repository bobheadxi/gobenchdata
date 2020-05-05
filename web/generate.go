package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//go:generate npm run build
//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

func defaultConfigPath(dir string) string { return path.Join(dir, "gobenchdata-web.json") }

// GenerateApp dumps the web app template into the provided directory
func GenerateApp(dir string) error {
	// clear directory of everything except config
	appConfigPath := defaultConfigPath(dir)
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot access target directory: %w", err)
	}
	for _, f := range dirs {
		fullName := path.Join(dir, f.Name())
		if f.IsDir() {
			os.RemoveAll(fullName)
		} else if fullName != appConfigPath {
			os.Remove(fullName)
		}
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
		if _, err := os.Stat(appConfigPath); os.IsExist(err) {
			return fmt.Errorf("found existing configuration: %w", err)
		}
	} else {
		os.Remove(appConfigPath)
	}

	// generate configuration
	b, _ := json.MarshalIndent(&defaultConfig, "", "  ")
	if err := ioutil.WriteFile(appConfigPath, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	return nil
}
