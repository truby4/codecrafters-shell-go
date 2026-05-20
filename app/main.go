package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		cmd, err := readCommand()
		if err != nil {
			fmt.Println(err.Error())
		}

		if cmd == "exit" {
			break
		} else if strings.HasPrefix(cmd, "echo ") {
			fmt.Println(cmd[5:])
		} else {
			fmt.Printf("%s: command not found\n", cmd)
		}
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
