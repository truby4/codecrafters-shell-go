package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type builtinFunc func([]string)

var builtins map[string]builtinFunc

func init() {
	builtins = map[string]builtinFunc{
		"echo": handleEcho,
		"type": handleType,
		"exit": handleExit,
		"pwd":  handlePWD,
		"cd":   handleCD,
	}
}

func main() {

	for {
		input, err := readInput()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		cmd, after := splitInput(input)

		args, err := lex(after)
		if err != nil {
			fmt.Printf("err: %s", err.Error())
			continue
		}

		if fn, exists := builtins[cmd]; exists {
			fn(args)
			continue
		}

		_, err = exec.LookPath(cmd)
		if err != nil {
			fmt.Printf("%s: command not found\n", cmd)
			continue
		}

		err = executePathCommand(cmd, args)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func executePathCommand(cmd string, args []string) error {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return fmt.Errorf("cmd execution error: %v", err)

	}
	fmt.Print(string(out))
	return nil
}

func readInput() (string, error) {
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading command: %v", err)
	}
	line = strings.TrimSuffix(line, "\n")
	return line, nil
}

func splitInput(input string) (cmd string, after string) {
	cmd = strings.Fields(input)[0]
	after = strings.TrimPrefix(input, cmd)
	return cmd, strings.TrimSpace(after)
}

func handleExit(_ []string) {
	os.Exit(0)
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleType(args []string) {
	for _, arg := range args {
		if _, exists := builtins[arg]; exists {
			fmt.Printf("%s is a shell builtin\n", arg)
			return
		}

		s, err := exec.LookPath(arg)
		if err != nil {
			fmt.Printf("%s: not found\n", arg)
			return
		}
		fmt.Printf("%s is %s\n", arg, s)
	}
}

func handlePWD(_ []string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
}

func handleCD(args []string) {
	if len(args) > 1 {
		fmt.Printf("cd: too many arguments")
		return
	}

	arg := args[0]

	if arg == "~" {
		arg = os.Getenv("HOME")
	}
	err := os.Chdir(arg)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", arg)
	}
}

func lex(after string) ([]string, error) {
	var tokens []string
	var current []rune

	inSingleQuote := false
	inDoubleQuote := false

	for _, ch := range after {
		switch ch {

		// toggle single quote tracker
		case '\'':
			if inDoubleQuote {
				// single quote inside double
				// quote is part of the token
				current = append(current, ch)
			} else {
				inSingleQuote = !inSingleQuote
			}

		case '"':
			inDoubleQuote = !inDoubleQuote

		case ' ', '\t', '\n':
			if inSingleQuote || inDoubleQuote {
				// spaces are part of the token
				current = append(current, ch)
			} else {
				// spaces outside quotes means to split the token
				if len(current) > 0 {
					tokens = append(tokens, string(current))
					current = nil
				}
			}

		default:
			current = append(current, ch)
		}
	}

	if inSingleQuote {
		return nil, fmt.Errorf("unterminated single quote")
	}

	if len(current) > 0 {
		tokens = append(tokens, string(current))
	}

	return tokens, nil
}
