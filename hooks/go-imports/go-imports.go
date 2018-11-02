package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/inturn/pre-commit-gobuild/internal/helpers"
)

func main() {
	if _, err := exec.Command("go", "get", "-u", "-v", "golang.org/x/tools/cmd/goimports").Output(); err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirs := helpers.DirsWith(wd, "\\.go$")

	errs := make([]error, 0)

	for _, d := range dirs {
		if !strings.Contains(d, "/vendor/") {
			relPath := strings.Replace(d, wd, ".", 1)
			cmd := exec.Command("goimports", "-w", "-l", relPath)
			if _, err := cmd.Output(); err != nil {
				log.Println(err)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 0 {
		log.Fatal(errs)
	}
}
