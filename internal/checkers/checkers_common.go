package checkers

import (
	"fmt"
)

type Piece struct {
	ID     int    `json:"id"`
	Color  string `json:"color"`
	IsKing bool   `json:"is_king"`
}

type Coords struct {
	Row int
	Col int
}

const pieceWhite = "w"
const pieceBlack = "b"

type Checkerscfg struct {
	Board           [8][8]Piece
	Pieces          map[int]Coords
	IsWhiteTurn     bool
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

var MovesMap = map[string]moveDir{
	"l":  moveLeft,
	"r":  moveRight,
	"bl": moveBackLeft,
	"br": moveBackRight,
}

type Move struct {
	Row       int
	Col       int
	Direction moveDir
}

// GetPlayerColor - returns pieceWhite if it is player one's turn, and pieceBlack if it is player 2's turn
func GetPlayerColor(isWhiteTurn bool) string {
	if isWhiteTurn {
		return pieceWhite
	}
	return pieceBlack
}

func (cfg *Checkerscfg) EndTurn() bool {
	cfg.IsWhiteTurn = !cfg.IsWhiteTurn
	if cfg.WhitePieceCount == 0 {
		fmt.Println("Black Wins!")
		return true
	} else if cfg.BlackPieceCount == 0 {
		fmt.Println("White Wins!")
		return true
	}

	if cfg.IsWhiteTurn {
		fmt.Println("White's Turn:")
	} else {
		fmt.Println("Black's Turn:")
	}

	return false
}

func isOutOfBounds(row, col int) bool {
	return row < 0 || row > 7 || col < 0 || col > 7
}

func applyDirection(row, col int, direction moveDir) (int, int) {
	switch direction {
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
	switch direction {
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

func (cfg *Checkerscfg) clearBoard() {
	cfg.Board = [8][8]Piece{}
}
func GetActualID(color string, id int) int {
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
