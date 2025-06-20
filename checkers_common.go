package main

import (
	"fmt"
)


type Piece struct {
	ID int
	Color string
	IsKing bool
}

type Coords struct {
	Row int
	Col int
}

const pieceWhite = "w"
const pieceBlack = "b"

type checkersCfg struct {
	Board [][]Piece
	Pieces map[int]Coords
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


// getPlayerColor - returns pieceWhite if it is player one's turn, and pieceBlack if it is player 2's turn
func (cfg *checkersCfg) getPlayerColor() string {
	if cfg.IsWhiteTurn {
		return pieceWhite
	} else {
		return pieceBlack
	}
}

func (cfg *checkersCfg) endTurn() bool {
	cfg.IsWhiteTurn = !cfg.IsWhiteTurn
	cfg.displayBoard()
	if cfg.WhitePieceCount == 0 {
		fmt.Println("Black Wins!")
		return true
	} else if cfg.BlackPieceCount == 0{
		fmt.Println("White Wins!")
		return true
	}

	if cfg.IsWhiteTurn{
		fmt.Println("White's Turn:")
	} else {
		fmt.Println("Black's Turn:")
	}

	return false
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
func getActualID(color string, id int) int {
	//don't rlly want to throw an error here, the -1 should at least tell me something has gone wrong
	if id == 0 {
		return -1
	}

	if color == pieceWhite {
		return id + 100
	} else {
		return id + 200
	}
}