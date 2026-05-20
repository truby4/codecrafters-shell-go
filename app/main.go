package main

import (
	"errors"
	"fmt"
)

func main() {
	var cmd string
	fmt.Print("$ ")
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		panic(errors.New("unable to scan command"))
	}

	fmt.Printf("%s: command not found\n", cmd)
}
