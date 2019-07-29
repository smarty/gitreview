package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func execute(dir, command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func prompt(message string) {
	log.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}
