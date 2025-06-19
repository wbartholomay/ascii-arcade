package main

import (
	"fmt"
	"strings"
)


type Piece struct {
	ID int
	Color string
	IsKing bool
}

const pieceWhite = "w"
const pieceBlack = "b"

type checkersCfg struct {
	Board [][]Piece
	IsWhiteTurn bool
	WhitePieceCount int
	BlackPieceCount int
}

type moveDir int

const (
	moveLeft moveDir = iota
	moveRight
	moveBackLeft
	moveBackRight
)

var movesMap = map[string]moveDir{
	"l":  moveLeft,
	"r":  moveRight,
	"bl": moveBackLeft,
	"br": moveBackRight,
}

type Move struct {
	Row int
	Col int
	Direction moveDir
}


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
		fmt.Println("   |       |       |       |       |       |       |       |       |")
		rowStr := fmt.Sprintf("%v  |", string(rune('a' + rowNum)))

		colNum := 0
		if !cfg.IsWhiteTurn {
			colNum = 7
		}

		for ; checkIndex(colNum); colNum += increment{
			piece := cfg.Board[rowNum][colNum]
			pieceStr := piece.Color
			if pieceStr == "" {
				pieceStr = " "
			}
			if piece.IsKing {
				pieceStr = strings.ToUpper(pieceStr)
			}
			rowStr += fmt.Sprintf("   %v   |", pieceStr)
		}
		fmt.Println(rowStr)
		fmt.Println("   |       |       |       |       |       |       |       |       |")
	}
	fmt.Println("   —————————————————————————————————————————————————————————————————")
}

// GetCurrentPieces - returns pieceWhite if it is player one's turn, and pieceBlack if it is player 2's turn
func (cfg *checkersCfg) getPlayerColor() string {
	if cfg.IsWhiteTurn {
		return pieceWhite
	} else {
		return pieceBlack
	}
}

func (cfg *checkersCfg) endTurn() error {
	cfg.IsWhiteTurn = !cfg.IsWhiteTurn
	cfg.displayBoard()
	fmt.Printf("White Pieces Remaining: %v    Black Pieces Remaining: %v\n", cfg.WhitePieceCount, cfg.BlackPieceCount)
	if cfg.WhitePieceCount == 0 {
		//TODO: Add logic to end game (will need to add bool return somewhere to be passed up the stack to the repl)
		fmt.Println("White Wins!")
	} else if cfg.BlackPieceCount == 0{
		fmt.Println("Black Wins!")
	}

	if cfg.IsWhiteTurn{
		fmt.Println("White's Turn:")
	} else {
		fmt.Println("Black's Turn:")
	}

	return nil
}

func isOutOfBounds(row, col int) bool {
	return row < 0 || row > 7 || col < 0 || col > 7
}

func applyDirection(row, col int, direction moveDir) (int, int){
	switch direction{
	case moveLeft:
		row -= 1
		col -= 1
	case moveRight:
		row -= 1
		col += 1
	case moveBackLeft:
		row += 1
		col -= 1
	case moveBackRight:
		row += 1
		col += 1
	}

	return row, col
}

func convertDirection(direction moveDir) moveDir {
	switch direction{
	case moveLeft:
		return moveBackRight
	case moveRight:
		return moveBackLeft
	case moveBackLeft:
		return moveRight
	case moveBackRight:
		return moveLeft
	default:
		return moveLeft
	}
}

func (cfg *checkersCfg) clearBoard() {
	cfg.Board = make([][]Piece, 8)
	for i := range cfg.Board {
		cfg.Board[i] = make([]Piece, 8)
	}
}