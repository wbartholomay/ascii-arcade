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
	if len(params) < 2{
		return errors.New("not enough arguments. Expecting move <piece number> <direction>")
	}

	pieceNum, err := strconv.ParseInt(params[0], 10, 8)
	if err != nil {
		return fmt.Errorf("expected a number, got: %v", params[0])
	}

	pieceId := getActualID(cfg.getPlayerColor(), int(pieceNum))
	piece, ok := cfg.Pieces[pieceId]
	if !ok {
		return fmt.Errorf("invalid piece number: %v", params[0])
	}

	directionString := params[1]
	direction, ok := movesMap[directionString]
	if !ok {
		return errors.New("invalid move direction. Valid moves are 'l', 'r', 'bl', 'br'")
	}

	
	move := Move{
		Row: piece.Row,
		Col: piece.Col,
		Direction: direction,
	}

	if err = cfg.movePiece(move); err != nil {
		return err
	}

	gameOver := cfg.endTurn()
	if gameOver {
		os.Exit(0)
	}
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
		capturedPiece = true
		if cfg.getPlayerColor() == pieceWhite {
			fmt.Println("Captured a black piece!")
			cfg.BlackPieceCount--
		} else {
			fmt.Println("Captured a white piece!")
			cfg.WhitePieceCount--
		}
	}

	tmpPiece := cfg.Pieces[piece.ID]
	tmpPiece.Row, tmpPiece.Col = targetRow, targetCol
	cfg.Pieces[piece.ID] = tmpPiece
	
	cfg.Board[targetRow][targetCol] = piece
	cfg.Board[move.Row][move.Col] = Piece{}

	//attempt king
	if piece.Color == pieceWhite && targetRow == 0 && !piece.IsKing{
		cfg.Board[targetRow][targetCol].IsKing = true
	} else if piece.Color == pieceBlack && targetRow == 7 && !piece.IsKing {
		cfg.Board[targetRow][targetCol].IsKing = true
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
		if isOutOfBounds(targetRow, targetCol) {
			continue
		}
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