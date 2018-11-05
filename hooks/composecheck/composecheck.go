package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	composeFiles := make([]string, 0)

	filepath.Walk(wd, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(`^docker-compose\\.yml$`, f.Name())
			if err == nil && r {
				composeFiles = append(composeFiles, path)
			}
		}
		return nil
	})

	errs := make([]error, 0)

	for _, fName := range composeFiles {
		f, err := os.Open(fName)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open file: %s --> %s", fName, err))
			continue
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read file: %s --> %s", fName, err))
			continue
		}

		dckComp := DockerCompose{}
		if err := yaml.Unmarshal(data, &dckComp); err != nil {
			errs = append(errs, fmt.Errorf("failed to unmarshal file: %s --> %s", fName, err))
			continue
		}

		// Validation
		if dckComp.Version != "3.7" {
			errs = append(errs, fmt.Errorf("required docker-compose version is 3.7: %s", f.Name()))
		}

		for _, svc := range dckComp.Services {
			if m, err := regexp.MatchString(`(?i)^\s*.*:latest\s*$`, svc.Image); err != nil || m {
				errs = append(errs, fmt.Errorf("image should not use the latest tag: %s", f.Name()))
			}
		}
	}

	if len(errs) != 0 {
		for _, e := range errs {
			log.Println(e)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

type DockerCompose struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Image string `yaml:"image"`
}
