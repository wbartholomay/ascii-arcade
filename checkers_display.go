package main

import (
	"fmt"
)

func (cfg *checkersCfg) displayBoard() {

	rowNum := 0
	increment := 1
	checkIndex := func (i int) bool {
		if cfg.IsWhiteTurn {
			return i < 8
		} else {
			return i >= 0
		}
	}

	if !cfg.IsWhiteTurn {
		rowNum = 7
		increment = -1
		fmt.Println("       7       6       5       4       3       2       1       0    ")
	} else {
		fmt.Println("       0       1       2       3       4       5       6       7    ")
	}

	for ; checkIndex(rowNum); rowNum += increment{
		fmt.Println("   —————————————————————————————————————————————————————————————————")
		squareStr := ""
		if (rowNum % 2 == 0 && cfg.IsWhiteTurn) || (rowNum % 2 != 0 && !cfg.IsWhiteTurn){
			squareStr = "   |       |#######|       |#######|       |#######|       |#######|"
		} else {
			squareStr = "   |#######|       |#######|       |#######|       |#######|       |"
		}
		fmt.Println(squareStr)
		rowStr := fmt.Sprintf("%v  |", string(rune('a' + rowNum)))

		colNum := 0
		if !cfg.IsWhiteTurn {
			colNum = 7
		}

		for ; checkIndex(colNum); colNum += increment{
			piece := cfg.Board[rowNum][colNum]
			if (rowNum % 2 == colNum % 2){
				rowStr += fmt.Sprintf("%v|", piece.renderPiece())
			} else {
				rowStr += "#######|"
			}
		}
		fmt.Println(rowStr)
		fmt.Println(squareStr)
	}
	fmt.Println("   —————————————————————————————————————————————————————————————————")
}

func (piece *Piece) renderPiece() string{
	if piece.Color == "" {
		return "       "
	}

	pieceStr := ""
	if piece.IsKing {
		pieceStr += "👑"
	} else {
		pieceStr += "  "
	}

	if piece.Color == pieceWhite {
		pieceStr += "⚪"
	} else if piece.Color == pieceBlack {
		pieceStr += "🔵"
	}
	pieceStr += toSubscript(piece.getDisplayID())

	if piece.getDisplayID() < 10 {
		pieceStr += " "
	}
	

	return pieceStr + " "
}

func toSubscript(n int) string {
	subs := []string{"", "₁", "₂", "₃", "₄", "₅", "₆", "₇", "₈", "₉", "₁₀", "₁₁", "₁₂"}
	return subs[n]
}

//This function kinda sucks but its temporary anyway
func (piece *Piece) getDisplayID() int {
	displayId := 0
	if piece.Color == pieceWhite {
		displayId = piece.ID - 100
	} else {
		displayId = piece.ID - 200
	}
	return displayId
}
