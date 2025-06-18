package checkers

import "fmt"

func displayBoard(cfg *checkersCfg, params ...string) error{
	fmt.Println("---------------------------------")
	for _, row := range cfg.Board{
		fmt.Println("---------------------------------")
		rowStr := "|"
		for _, square := range row{
			rowStr += fmt.Sprintf(" %v |", square.pieceStatus)
		}
		fmt.Println(rowStr)
	}
	fmt.Println("---------------------------------")

	return nil
}