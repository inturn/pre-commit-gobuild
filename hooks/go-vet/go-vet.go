package main

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	verStr := strings.TrimPrefix(runtime.Version(), "go")

	if strings.Compare("1.12", verStr) == 1 {
		log.Fatalf("Update your go version. %s --> 1.12+", verStr)
	}

	cmd := exec.Command("go", "vet", "./...")
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(res))
	}
}
