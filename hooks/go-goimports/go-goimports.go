package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/inturn/pre-commit-gobuild/internal/helpers"
)

func main() {
	if _, err := exec.Command("go", "get", "-u", "-v", "golang.org/x/tools/cmd/goimports").Output(); err != nil {
		log.Fatal(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirs := helpers.DirsWith(workDir, "\\.go$")

	wg := &sync.WaitGroup{}

	for _, dir := range dirs {
		go func(d, wd string) {
			wg.Add(1)
			defer wg.Done()
			if !strings.Contains(d, "/vendor/") {
				relPath := strings.Replace(d, wd, ".", 1)
				cmd := exec.Command("goimports", "-w", "-l", relPath)
				if _, err := cmd.Output(); err != nil {
					log.Printf("error occured on goimports execute: %s\n", err)
				}
			}
		}(dir, workDir)
	}
	wg.Wait()
}
