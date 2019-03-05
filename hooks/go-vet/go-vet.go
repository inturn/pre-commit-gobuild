package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	verStr := strings.TrimPrefix(runtime.Version(), "go")

	ver, err := strconv.ParseFloat(verStr, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var arg []string

	if ver < 1.12 {
		arg = []string{"tool", "vet", "./..."}
	} else {
		arg = []string{"vet", "./..."}
	}

	cmd := exec.Command("go", arg...)
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(res))
		os.Exit(1)
	}
}
