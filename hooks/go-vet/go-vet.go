package main

import (
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	verStr := strings.TrimPrefix(runtime.Version(), "go")

	ver, err := strconv.ParseFloat(verStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	if ver < 1.12 {
		log.Fatalf("Update your go version. %f --> 1.12+", ver)
	}

	cmd := exec.Command("go", "vet", "./...")
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(res))
	}
}
