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

	row := int8(rowRune - 'a')

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
	

	if err = movePiece(cfg, int8(row), int8(col), direction); err != nil {
		return err
	}

	cfg.displayBoard()
	cfg.endTurn()
	return nil
}