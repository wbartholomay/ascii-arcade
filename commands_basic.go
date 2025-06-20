package main

import (
	"fmt"
	"os"
)

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

func commandBoard(cfg *checkersCfg, params ...string) error{
	cfg.displayBoard()
	return nil
}

func commandConcede(cfg *checkersCfg, params ...string) error{
	if cfg.IsWhiteTurn {
		fmt.Println("White conceded, black wins!")
	} else {
		fmt.Println("Black conceded, white wins!")
	}

	os.Exit(0)
	return nil
}