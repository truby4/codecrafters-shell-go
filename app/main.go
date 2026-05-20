package main

import (
	"bufio"
	"fmt"
	"os"
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
			if slices.Contains(shellBuiltins, after) {
				fmt.Printf("%s is a shell builtin\n", after)
			} else {
				fmt.Printf("%s: not found\n", after)
			}
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
