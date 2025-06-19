package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
	direction, ok := movesMap[directionString]
	if !ok {
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
	piece := cfg.Board[move.Row][move.Col]

	//check selected square
	if piece.Color == "" {
		return errors.New("no piece on this square")
	}
	if piece.Color != cfg.getPlayerColor() {
		return errors.New("you can only move your own pieces")
	}
	if (move.Direction == moveBackLeft || move.Direction == moveBackRight) && !piece.IsKing {
		return errors.New("only kings can move backwards")
	}

	//get absolute direction based on the input direction and the piece color
	if !cfg.IsWhiteTurn {
		move.Direction = convertDirection(move.Direction)
	}
	
	//validate move
	targetRow, targetCol := applyDirection(move.Row, move.Col, move.Direction)
	capturedPiece := false

	if isOutOfBounds(targetRow, targetCol){
		return errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[targetRow][targetCol].Color == cfg.getPlayerColor() {
		return errors.New("target square is occupied")
	}
	if cfg.Board[targetRow][targetCol].Color != cfg.getPlayerColor() && cfg.Board[targetRow][targetCol].Color != ""{
		//attempt capture
		captureRow, captureCol := targetRow, targetCol
		targetRow, targetCol = applyDirection(targetRow, targetCol, move.Direction)
		if isOutOfBounds(targetRow, targetCol){
			return errors.New("cannot move a piece outside of the board")
		}

		if cfg.Board[targetRow][targetCol].Color != "" {
			return errors.New("target square is occupied")
		}

		cfg.Board[captureRow][captureCol] = Piece{}
		if cfg.getPlayerColor() == pieceWhite {
			cfg.BlackPieceCount--
		} else {
			cfg.WhitePieceCount--
		}
	}

	cfg.Board[targetRow][targetCol] = piece
	cfg.Board[move.Row][move.Col] = Piece{}

	//attempt king
	if piece.Color == pieceWhite && targetRow == 0 && !piece.IsKing{
		piece.IsKing = true
	} else if piece.Color == pieceBlack && targetRow == 7 && !piece.IsKing {
		piece.IsKing = true
	}
	if !capturedPiece {
		return nil
	}

	//check for double capture
	if capturedPiece {
		nextMoves := cfg.checkSurroundingSquaresForCapture(targetRow, targetCol)
		if len(nextMoves) == 0 {
			return nil
		}

		cfg.displayBoard()
		fmt.Print("Another capture is available, enter one of the following directions: ")
		for _, moveStr := range nextMoves {
			fmt.Printf("%v, ", moveStr)
		}
		fmt.Println()
		input := ""
		for {
			fmt.Print("Checkers > ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input = strings.ToLower(scanner.Text())
			if !slices.Contains(nextMoves, input){
				fmt.Println("Please enter one of the displayed directions.")
			} else {
				break
			}
		}
		return cfg.movePiece(Move{
			Row: targetRow,
			Col: targetCol,
			Direction: movesMap[input],
		})
	}

	return nil
}

func (cfg *checkersCfg) checkSurroundingSquaresForCapture(row, col int) []string{
	piece := cfg.Board[row][col]
	captureMoves := []string{}

	//if the piece is not a king, only check forward moves. Otherwise, check all directions
	moves := []string{"l", "r"}
	if piece.IsKing {
			moves = append(moves, "bl", "br")
	}
	
	for _, moveStr := range moves {
		move := movesMap[moveStr]

		if !cfg.IsWhiteTurn {
			move = convertDirection(move)
		}

		targetRow, targetCol := applyDirection(row, col, move)
		fmt.Printf("Target row: %v   Target col: %v   Move direction: %v\n", targetRow, targetCol, moveStr)

		if cfg.Board[targetRow][targetCol].Color == cfg.getPlayerColor() || cfg.Board[targetRow][targetCol].Color == "" {
			continue
		}
		targetRow, targetCol = applyDirection(targetRow, targetCol, move)
		if isAvailable, _ := cfg.isSquareAvailable(targetRow, targetCol); isAvailable {
			captureMoves = append(captureMoves, moveStr)
		}
	}

	return captureMoves
}

func (cfg *checkersCfg) isSquareAvailable(row, col int) (bool, error){
	if isOutOfBounds(row, col){
			return false, errors.New("cannot move a piece outside of the board")
		}

	if cfg.Board[row][col].Color != "" {
		return false, errors.New("target square is occupied")
	}

	return true, nil
}