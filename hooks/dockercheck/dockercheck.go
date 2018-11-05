package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/inturn/pre-commit-gobuild/internal/dockerfile"
	"github.com/inturn/pre-commit-gobuild/internal/dockerfile/command"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dckFiles := make([]string, 0)

	filepath.Walk(wd, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(`^Dockerfile.*$`, f.Name())
			if err == nil && r {
				dckFiles = append(dckFiles, path)
			}
		}
		return nil
	})

	errs := make([]error, 0)
	defer func() {
		if len(errs) != 0 {
			for _, e := range errs {
				log.Println(e)
			}
			os.Exit(1)
		}
		os.Exit(0)
	}()

	for _, f := range dckFiles {
		directive := parser.Directive{}
		parser.SetEscapeToken("`", &directive)

		f, err := os.Open(f)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		node, err := parser.Parse(f, &directive)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, n := range node.Children {
			v := strings.ToUpper(n.Value)
			if m, err := regexp.MatchString(fmt.Sprintf(`^%s.*$`, v), n.Original); err != nil || !m {
				errs = append(errs, fmt.Errorf("capitalize Dockerfile Instructions: %s line %d", f.Name(), n.StartLine))
			}
			if n.Value == command.From {
				if m, err := regexp.MatchString(`(?i)^(from)\s.*:latest\s*$`, n.Original); err != nil || m {
					errs = append(errs, fmt.Errorf("images should not use the latest tag: %s line %d", f.Name(), n.StartLine))
				}
			}
			if m, err := regexp.MatchString(`^.*apt-get.*install.*$`, n.Original); err != nil || m {
				if !strings.Contains(n.Original, "--no-install-recommends") {
					fmt.Printf("Consider `--no-install-recommends`: %s line %d", f.Name(), n.StartLine)
				}
			}
		}
	}
}
