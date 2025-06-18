package main

import "fmt"

func commandHelp(cfg *checkersCfg, params ...string) error{
	fmt.Print("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandExit(cfg *checkersCfg, params ...string) error {
	fmt.Println("Closing ASCII Checkers... Goodbye!")

	return nil
}