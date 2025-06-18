package main

type tileColor int8
type pieceStatus int8

const (
	tileWhite tileColor = iota
	tileBlack
)

const (
	pieceEmpty pieceStatus = iota
	pieceWhite
	pieceBlack
)

type tileStatus struct{
	tileColor tileColor
	pieceStatus pieceStatus
}

type checkersCfg struct {
	Board [][]tileStatus
	PlayerPiece pieceStatus
}

func startCheckers() checkersCfg {
	//initialize board
	board := initializeBoard()

	return checkersCfg{
		Board: board,
		PlayerPiece: pieceWhite,
	}
}

func initializeBoard() [][]tileStatus {
	board := make([][]tileStatus, 8)
	for row := range board {
		board[row] = make([]tileStatus, 8)
		for col := range board[row] {
			board[row][col] = tileStatus{}

			//initialize tile colors
			isRowEven := (row % 2) == 0
			isColumnEven := (col % 2) == 0
			if isRowEven == isColumnEven {
				board[row][col].tileColor = tileBlack
			} else {
				board[row][col].tileColor = tileWhite
			}

			//initialize pieces
			if board[row][col].tileColor == tileBlack && row < 3 {
				board[row][col].pieceStatus = pieceBlack
			} else if board[row][col].tileColor == tileBlack && row > 4 {
				board[row][col].pieceStatus = pieceWhite
			} else {
				board[row][col].pieceStatus = pieceEmpty
			}
		}
	}

	return board
}