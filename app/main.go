package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var shellBuiltins = []string{
	"echo", "type", "exit", "pwd", "cd",
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

		if cmd == "pwd" {
			handlePWD()
			continue
		}

		after, found := strings.CutPrefix(cmd, "cd")
		if found {
			handleCD(strings.TrimSpace(after))
			continue
		}

		after, found = strings.CutPrefix(cmd, "echo ")
		if found {
			fmt.Println(after)
			continue
		}

		after, found = strings.CutPrefix(cmd, "type ")
		if found {
			handleType(after)
			continue
		}

		split_cmd := strings.Fields(cmd)
		_, err = exec.LookPath(split_cmd[0])
		if err != nil {
			fmt.Printf("%s: command not found\n", split_cmd[0])
			continue
		}

		out, err := exec.Command(split_cmd[0], split_cmd[1:]...).Output()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Print(string(out))
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

func handlePWD() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
}

func handleCD(arg string) {
	if arg == "~" {
		arg = os.Getenv("HOME")
	}
	err := os.Chdir(arg)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", arg)
	}
}
