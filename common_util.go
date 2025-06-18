package main

import (
	"fmt"
)


func (cfg *checkersCfg) displayBoard() error {

	rowNum := 0
	increment := 1
	checkIndex := func (i int) bool {
		if cfg.IsPlayerOneTurn {
			return i < 8
		} else {
			return i >= 0
		}
	}

	if !cfg.IsPlayerOneTurn {
		rowNum = 7
		increment = -1
	}

	fmt.Println("       0       1       2       3       4       5       6       7    ")
	for ; checkIndex(rowNum); rowNum += increment{
		fmt.Println("   —————————————————————————————————————————————————————————————————")
		fmt.Println("   |       |       |       |       |       |       |       |       |")
		rowStr := fmt.Sprintf("%v  |", string(rune('a' + rowNum)))

		colNum := 0
		if !cfg.IsPlayerOneTurn {
			colNum = 7
		}

		for ; checkIndex(colNum); colNum += increment{
			piece := cfg.Board[rowNum][colNum]
			pieceStr := piece.Color
			if pieceStr == "" {
				pieceStr = " "
			}
			rowStr += fmt.Sprintf("   %v   |", pieceStr)
		}
		fmt.Println(rowStr)
		fmt.Println("   |       |       |       |       |       |       |       |       |")
	}
	fmt.Println("   —————————————————————————————————————————————————————————————————")

	return nil
}

// GetCurrentPieces - returns pieceWhite if it is player one's turn, and pieceBlack if it is player 2's turn
func (cfg *checkersCfg) getCurrentPieces() string {
	if cfg.IsPlayerOneTurn {
		return pieceWhite
	} else {
		return pieceBlack
	}
}

func (cfg *checkersCfg) endTurn() error {
	if cfg.IsPlayerOneTurn{
		fmt.Println("Player 2's Turn:")
	} else {
		fmt.Println("Player 1's Turn:")
	}

	cfg.IsPlayerOneTurn = !cfg.IsPlayerOneTurn
	cfg.displayBoard()

	return nil
}