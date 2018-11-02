package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/inturn/pre-commit-gobuild/internal/helpers"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirs := helpers.DirsWith(wd, "\\.go$")

	for _, d := range dirs {
		if !strings.Contains(d, "/vendor/") {
			relPath := strings.Replace(d, wd, ".", 1)
			cmd := exec.Command("go", "install", relPath)
			out, err := cmd.CombinedOutput()
			outStr := string(out)
			if err != nil {
				if !strings.Contains(outStr, "no non-test Go files") && !strings.Contains(outStr, "build constraints exclude all Go files") {
					fmt.Println(fmt.Sprintf("building %s finished with error %s", relPath, err.Error()))
					fmt.Printf(string(out))
					os.Exit(1)
				}
			} else {
				fmt.Println(fmt.Sprintf("building %s finished ok", relPath))
			}
		}
	}
}
