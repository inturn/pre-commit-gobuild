package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/inturn/pre-commit-gobuild/internal/helpers"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirs := helpers.DirsWith(workDir, "\\.go$")

	wg := &sync.WaitGroup{}

	for _, dir := range dirs {
		wg.Add(1)
		go func(d, wd string) {
			defer wg.Done()
			if !strings.Contains(d, "/vendor/") {
				files, err := ioutil.ReadDir(d)
				if err != nil {
					log.Printf("error occured on read dir %s: %s", d, err)
				}
				for _, f := range files {
					if !strings.HasSuffix(f.Name(), ".go") {
						continue
					}
					cmd := exec.Command("goimports", "-w", "-l", fmt.Sprintf("%s/%s", d, f.Name()))
					if _, err := cmd.Output(); err != nil {
						log.Printf("error occured on goimports execute: %s\n", err)
					}
				}
			}
		}(dir, workDir)
	}
	wg.Wait()
}
