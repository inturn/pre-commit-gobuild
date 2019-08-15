package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "generate", "./...")
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(res))
	}
}
