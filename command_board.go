package main

import "fmt"

func displayBoard(cfg *checkersCfg, params ...string) error{
	fmt.Println("      0     1     2     3     4     5     6     7   ")
	for i, row := range cfg.Board{
		fmt.Println("   —————————————————————————————————————————————————")
		rowStr := fmt.Sprintf("%v  |", string(rune('a' + i)))
		for _, square := range row{
			rowStr += fmt.Sprintf("  %v  |", square.pieceStatus)
		}
		fmt.Println(rowStr)
	}
	fmt.Println("   —————————————————————————————————————————————————")

	return nil
}