package checkers

import "errors"

type moveDir int8

const (
	moveLeft moveDir = iota
	moveRight
	moveBackLeft
	moveBackRight
)

// movePiece - takes the initial and direction to move piece
// validates the move can be made, and if it can the board is updated
func movePiece(cfg *checkersCfg, startRow, startCol int8, direction moveDir) error{
	if direction > 3 {
		return errors.New("invalid move option")
	}

	destRow, destCol := startRow, startCol
	piece := cfg.Board[startRow][startCol].pieceStatus

	//TODO: update this to check if the piece is the players
	if piece == pieceEmpty {
		return errors.New("no piece on this square")
	}

	//get absolute direction based on the inputted direction and the piece color
	absoluteDir := direction
	if piece == pieceBlack {
		absoluteDir = convertDirection(absoluteDir)
	}

	//validate move
	switch absoluteDir{
	case moveLeft:
		destRow += 2
		destCol -= 2
	case moveRight:
		destRow += 2
		destCol += 2
	case moveBackLeft:
		destRow -= 2
		destCol -= 2
	case moveBackRight:
		destRow -= 2
		destCol += 2
	}
	if err := validateMove(*cfg, destRow, destCol); err != nil {
		return err
	}

	//TODO: add capture logic?
	//update board
	cfg.Board[destRow][destCol].pieceStatus = piece
	cfg.Board[startRow][startCol].pieceStatus = pieceEmpty
	return nil
}

func validateMove(cfg checkersCfg, row, col int8) error {

	if row < 0 || row > 7 || col < 0 || col > 7 {
		return errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[row][col].pieceStatus != pieceEmpty {
		return errors.New("there is already a piece on that square")
	}

	//TODO: add logic to check if a piece is a king, and allow backwards moves only if it is a king
	return nil
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