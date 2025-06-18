package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
)


type cliCommand struct {
	name string
	description string
	callback func(params ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"help" : {
			name: "help",
			description: "Displays a list of commmands",
			callback: commandHelp,
		},
		"exit" : {
			name: "exit",
			description: "Exits the application",
			callback: commandExit,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	substrings := strings.Fields(text)
	return substrings
}

//TODO: wrap this in an outer repl if more games than checkers exit
func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to ASCII Checkers!\nEnter 'help' to see a list of commands.")

	for {
		fmt.Print("Checkers > ")
		scanner.Scan()

		t := scanner.Text()
		input := cleanInput(t)
		if len(input) == 0 { continue }

		cmd, ok := getCommands()[input[0]]
		if !ok{
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(input[1:]...)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}