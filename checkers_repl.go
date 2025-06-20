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
	callback func(cfg *checkersCfg, params ...string) error
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
			description: "Exits checkers",
			callback: commandExit,
		},
		"move" : {
			name: "move",
			description: "Move a piece. Takes arguments <row> <col> <direction {'l', 'r', 'bl', 'br'}>",
			callback: commandMove,
		},
		"board" : {
			name: "board",
			description: "Displays the current board state",
			callback: commandBoard,
		},
		"concede" : {
			name: "concede",
			description: "Concede the game",
			callback: commandConcede,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	substrings := strings.Fields(text)
	return substrings
}

//TODO: wrap this in an outer repl if more games than checkers exit
func StartCheckersRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to ASCII Checkers!")

	cfg := startCheckers()
	cfg.displayBoard()
	fmt.Println("Player 1's Turn:")

	for {
		fmt.Print("Checkers > ")
		scanner.Scan()

		t := scanner.Text()
		input := cleanInput(t)
		if len(input) == 0 { 
			continue 
		}

		cmd, ok := getCommands()[input[0]]
		if !ok{
			fmt.Println("Unknown command. Enter 'help' to see a list of commands.")
			continue
		}

		//exit checkers game
		if cmd.name == "exit" {
			break
		}

		err := cmd.callback(&cfg, input[1:]...)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}