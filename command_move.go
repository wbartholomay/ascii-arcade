package main

import (
	"errors"
	"fmt"
	"strconv"
)

func commandMove(cfg *checkersCfg, params ...string) error {
	//validate params - expecting move <row> <col> <direction>
	if len(params) < 3{
		return errors.New("not enough arguments. Expecting move <row> <col> <direction>")
	}

	rowRune := params[0][0]

	if len(params[0]) > 1 || rowRune < 'a' || rowRune > 'f'{
		return errors.New("error parsing row arg. expecting 1 character between a and f")
	}

	row := int(rowRune - 'a')

	col, err := strconv.ParseInt(params[1], 10, 8)
	if err != nil {
		return fmt.Errorf("error parsing col arg to int: %w", err)
	}

	if row < 0 || row > 7 || col < 0 || col > 7 {
		return errors.New("row and col must be within range [0,7]")
	}

	directionString := params[2]
	var direction moveDir
	switch directionString{
	case "l":
		direction = moveLeft
	case "r":
		direction = moveRight
	case "bl":
		direction = moveBackLeft
	case "br":
		direction = moveBackRight
	default:
		return errors.New("invalid move direction. Valid moves are 'l', 'r', 'bl', 'br'")
	}
	
	move := Move{
		Row: row,
		Col: int(col),
		Direction: direction,
	}

	if err = cfg.movePiece(move); err != nil {
		return err
	}

	cfg.endTurn()
	return nil
}


// movePiece - takes the initial and direction to move piece
// validates the move can be made, and if it can the board is updated
func (cfg *checkersCfg) movePiece(move Move) error{
	if move.Direction > 3 {
		return errors.New("invalid move option")
	}

	destRow, destCol := move.Row, move.Col
	piece := cfg.Board[move.Row][move.Col]

	if piece.Color == "" {
		return errors.New("no piece on this square")
	}

	if piece.Color != cfg.getCurrentPieces() {
		return errors.New("you can only move your own pieces")
	}

	//get absolute direction based on the inputted direction and the piece color
	absoluteDir := move.Direction
	if cfg.getCurrentPieces() == pieceBlack {
		absoluteDir = convertDirection(absoluteDir)
	}
	
	//validate move
	switch absoluteDir{
	case moveLeft:
		destRow -= 1
		destCol -= 1
	case moveRight:
		destRow -= 1
		destCol += 1
	case moveBackLeft:
		destRow += 1
		destCol -= 1
	case moveBackRight:
		destRow += 1
		destCol += 1
	}

	if err := validateMove(*cfg, destRow, destCol); err != nil {
		return err
	}

	//update board
	cfg.Board[destRow][destCol] = piece
	cfg.Board[move.Row][move.Col] = Piece{}
	return nil
}

func validateMove(cfg checkersCfg, row, col int) error {

	if row < 0 || row > 7 || col < 0 || col > 7 {
		return errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[row][col].Color != "" {
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