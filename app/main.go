package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var shellBuiltins = []string{
	"echo", "type", "exit",
}

func main() {

	for {
		cmd, err := readCommand()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if cmd == "exit" {
			break
		}

		after, found := strings.CutPrefix(cmd, "echo ")
		if found {
			fmt.Println(after)
			continue
		}

		after, found = strings.CutPrefix(cmd, "type ")
		if found {
			handleType(after)
			continue
		}

		fmt.Printf("%s: command not found\n", cmd)
	}
}

func readCommand() (string, error) {
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading command: %v", err)
	}
	line = strings.TrimSuffix(line, "\n")
	return line, nil
}

func handleType(after string) {
	if slices.Contains(shellBuiltins, after) {
		fmt.Printf("%s is a shell builtin\n", after)
		return
	}

	s, err := exec.LookPath(after)
	if err != nil {
		fmt.Printf("%s: not found\n", after)
		return
	}
	fmt.Printf("%s is %s\n", after, s)
}
