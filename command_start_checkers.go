package main

import (
	"fmt"

	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

func commandStartCheckers(params ...string) error {
	cfg := checkers.StartCheckers()
	for row := range cfg.Board {
		fmt.Println(cfg.Board[row])
	}

	return nil
}