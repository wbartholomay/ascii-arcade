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
	if len(params[0]) > 1 || rowRune < 'a' || rowRune > 'h'{
		return errors.New("error parsing row arg. expecting 1 character between a and h")
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
		DestRow: row,
		DestCol: int(col),
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
	piece := cfg.Board[move.Row][move.Col]

	if piece.Color == "" {
		return errors.New("no piece on this square")
	}

	if piece.Color != cfg.getPlayerColor() {
		return errors.New("you can only move your own pieces")
	}

	//get absolute direction based on the inputted direction and the piece color
	if cfg.getPlayerColor() == pieceBlack {
		move.Direction = convertDirection(move.Direction)
	}
	
	//validate move
	move.applyDirection()

	if err := validateMove(cfg, &move); err != nil {
		return err
	}

	//update board
	cfg.Board[move.DestRow][move.DestCol] = piece
	cfg.Board[move.Row][move.Col] = Piece{}
	return nil
}

func validateMove(cfg *checkersCfg, move *Move) error {

	if isOutOfBounds(move.DestRow, move.DestCol){
		return errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[move.DestRow][move.DestCol].Color != "" {
		return attemptCapture(cfg, move)
	}

	//TODO: add logic to check if a piece is a king, and allow backwards moves only if it is a king
	return nil
}

func attemptCapture(cfg *checkersCfg, move *Move) error {
	if cfg.Board[move.DestRow][move.DestCol].Color == cfg.getPlayerColor() {
		return errors.New("there is already a piece on that square")
	}
	//check next tile
	captureRow, captureCol := move.DestRow, move.DestCol
	move.applyDirection()
	//if space behind piece is not open
	if isOutOfBounds(move.DestRow, move.DestCol){
		return errors.New("cannot move a piece outside of the board")
	}

	if !cfg.isTileEmpty(move.DestRow, move.DestCol) {
		return fmt.Errorf("there is already a piece on square (%v, %v)", move.DestRow, move.DestCol)
	}

	//capture piece
	if cfg.getPlayerColor() == "W" {
		cfg.BlackPieceCount--
	} else {
		cfg.WhitePieceCount--
	}
	cfg.Board[captureRow][captureCol] = Piece{}

	return nil
}