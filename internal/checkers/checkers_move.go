package checkers

import (
	"errors"
	"fmt"
)

// MovePiece - takes the initial and direction to move piece
// validates the move can be made, and if it can the board is updated
// takes in argument move and either a channel or a connection
func (cfg *checkersCfg) MovePiece(move Move, 
	transport Transport[ServerToClientData, ClientToServerData]) error {
	piece := cfg.Board[move.Row][move.Col]
	playerColor := GetPlayerColor(cfg.IsWhiteTurn)

	//check selected square
	if piece.Color == "" {
		return errors.New("no piece on this square")
	}
	if piece.Color != playerColor {
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

	if isOutOfBounds(targetRow, targetCol) {
		return errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[targetRow][targetCol].Color == playerColor {
		return errors.New("target square is occupied")
	}
	if cfg.Board[targetRow][targetCol].Color != playerColor && cfg.Board[targetRow][targetCol].Color != "" {
		//attempt capture
		captureRow, captureCol := targetRow, targetCol
		targetRow, targetCol = applyDirection(targetRow, targetCol, move.Direction)
		if isOutOfBounds(targetRow, targetCol) {
			return errors.New("cannot move a piece outside of the board")
		}

		if cfg.Board[targetRow][targetCol].Color != "" {
			return errors.New("target square is occupied")
		}

		delete(cfg.Pieces, cfg.Board[captureRow][captureCol].ID)
		cfg.Board[captureRow][captureCol] = Piece{}
		capturedPiece = true
		if playerColor == pieceWhite {
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
	if piece.Color == pieceWhite && targetRow == 0 && !piece.IsKing {
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

		DisplayBoard(cfg.Board, cfg.IsWhiteTurn)
		fmt.Print("Another capture is available, enter one of the following directions: ")
		for _, moveStr := range nextMoves {
			fmt.Printf("%v, ", moveStr)
		}
		fmt.Println()
		err := transport.SendData(ServerToClientData{
			Board:             cfg.Board,
			Pieces:            cfg.Pieces,
			IsDoubleJump:      true,
			DoubleJumpOptions: nextMoves,
		})
		if err != nil {
			return err
		}

		dataFromClient, err := transport.ReceiveData()
		if err != nil {
			return err
		}

		return cfg.MovePiece(Move{
			Row:       targetRow,
			Col:       targetCol,
			Direction: MovesMap[dataFromClient.DoubleJumpDirection],
		}, transport)
	}

	return nil
}

func (cfg *checkersCfg) checkSurroundingSquaresForCapture(row, col int) []string {
	piece := cfg.Board[row][col]
	captureMoves := []string{}

	//if the piece is not a king, only check forward moves. Otherwise, check all directions
	moves := []string{"l", "r"}
	if piece.IsKing {
		moves = append(moves, "bl", "br")
	}

	for _, moveStr := range moves {
		move := MovesMap[moveStr]

		if !cfg.IsWhiteTurn {
			move = convertDirection(move)
		}

		targetRow, targetCol := applyDirection(row, col, move)
		if isOutOfBounds(targetRow, targetCol) {
			continue
		}
		if cfg.Board[targetRow][targetCol].Color == GetPlayerColor(cfg.IsWhiteTurn) || cfg.Board[targetRow][targetCol].Color == "" {
			continue
		}
		targetRow, targetCol = applyDirection(targetRow, targetCol, move)
		if isAvailable, _ := cfg.isSquareAvailable(targetRow, targetCol); isAvailable {
			captureMoves = append(captureMoves, moveStr)
		}
	}

	return captureMoves
}

func (cfg *checkersCfg) isSquareAvailable(row, col int) (bool, error) {
	if isOutOfBounds(row, col) {
		return false, errors.New("cannot move a piece outside of the board")
	}

	if cfg.Board[row][col].Color != "" {
		return false, errors.New("target square is occupied")
	}

	return true, nil
}
