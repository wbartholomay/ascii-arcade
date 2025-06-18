package main

import (
	"fmt"
	"os"
)

func commandHelp(params ...string) error{
	fmt.Print("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandExit(params ...string) error {
	fmt.Println("Closing ASCII Checkers... Goodbye!")
	os.Exit(0)

	return nil
}