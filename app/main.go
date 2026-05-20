package main

import (
	"fmt"
)

func main() {
	for {
		var cmd string
		err := readCommand(&cmd)
		if err != nil {
			fmt.Println(err.Error())
		}

		if cmd == "exit" {
			break
		}

		fmt.Printf("%s: command not found\n", cmd)
	}
}

func readCommand(cmd *string) error {
	fmt.Print("$ ")
	_, err := fmt.Scanln(cmd)
	if err != nil {
		return fmt.Errorf("error reading command: %v", err)
	}
	return nil
}
