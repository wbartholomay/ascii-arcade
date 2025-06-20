package main

import (
	"reflect"
	"testing"
)

func TestMovePiece(t *testing.T) {
	tests := []struct {
		name string
		row int
		col int
		expectedRow int
		expectedCol int
		direction moveDir
		wantErr bool
	}{
		{
			name: "Successful move left",
			row: 5,
			col: 1,
			expectedRow: 4,
			expectedCol: 0,
			direction: moveLeft,
			wantErr: false,
		},
		{
			name: "Successful move right",
			row: 5,
			col: 1,
			expectedRow: 4,
			expectedCol: 2,
			direction: moveRight,
			wantErr: false,
		},
		{
			name: "Move out of bounds",
			row: 5,
			col: 7,
			direction: moveRight,
			wantErr: true,
		},
		{
			name: "Move into occupied space",
			row: 6,
			col: 0,
			direction: moveRight,
			wantErr: true,
		},
		{
			name: "Attempt to move backwards as not king",
			row: 6,
			col: 0,
			direction: moveBackLeft,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			cfg := startCheckers()
			err := cfg.movePiece(Move{
				Row: tt.row,
				Col: tt.col,
				Direction: tt.direction,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("movePiece() error = %v, wantErr %v. Starting square = (%v, %v), move direction = %v", err, tt.wantErr, tt.row, tt.col, tt.direction)
				return
			}

			if !tt.wantErr && cfg.Board[tt.expectedRow][tt.expectedCol].Color == ""{
				t.Errorf("movePiece() did not move piece to expected square: (%v, %v)", tt.expectedRow, tt.expectedCol)
				return
			}
		})
	}
}

func TestCapture(t *testing.T) {
	tests := []struct {
		name string
		row int
		col int
		expectedRow int
		expectedCol int
		direction moveDir
		wantErr bool
	}{
		{
			name: "Successful single capture",
			row: 4,
			col: 4,
			expectedRow: 2,
			expectedCol: 2,
			direction: moveLeft,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T){
			cfg := startCheckers()
			cfg.Board[4][4] = Piece{
				Color: pieceWhite,
			}
			cfg.Board[3][3] = Piece{
				Color: pieceBlack,
			}
			cfg.Board[2][2] = Piece{}
			err := cfg.movePiece(Move{
				Row: tt.row,
				Col: tt.col,
				Direction: tt.direction,
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("movePiece() error = %v, wantErr %v. Starting square = (%v, %v), move direction = %v", err, tt.wantErr, tt.row, tt.col, tt.direction)
			}

			if !tt.wantErr && cfg.Board[tt.expectedRow][tt.expectedCol].Color == ""{
				t.Errorf("movePiece() did not move piece to expected square: (%v, %v)", tt.expectedRow, tt.expectedCol)
				return
			}

			if !tt.wantErr && cfg.BlackPieceCount != 11 {
				t.Errorf("Capture unsuccessful, piece count did not decrease. BlackPieceCount: %v", cfg.BlackPieceCount)
			}
		})
	}
}
func TestCheckSurroundingSquaresForCapture(t *testing.T) {
    cfg := startCheckers()
	type testPiece struct {
		row int
		col int
		color string
		isKing bool
	}

	tests := []struct {
		name string
		row int
		col int
		pieces []testPiece
		isWhiteTurn bool
		expectedMoves []string
	}{
		{
			name: "No captures",
			row: 3,
			col: 3,
			isWhiteTurn: true,
			pieces: []testPiece{
				{
					row: 3,
					col: 3,
					color: pieceWhite,
				},
			},
			expectedMoves: []string{},
		},
		{
			name: "Capture to left",
			row: 3,
			col: 3,
			isWhiteTurn: true,
			pieces: []testPiece{
				{
					row: 3,
					col: 3,
					color: pieceWhite,
				},
				{
					row: 2,
					col: 2,
					color: pieceBlack,
				},
			},
			expectedMoves: []string{"l"},
		},{
			name: "Black king capture to right and back right",
			row: 3,
			col: 3,
			isWhiteTurn: false,
			pieces: []testPiece{
				{
					row: 3,
					col: 3,
					color: pieceBlack,
					isKing: true,
				},
				{
					row: 4,
					col: 2,
					color: pieceWhite,
				},
				{
					row: 2,
					col: 2,
					color: pieceWhite,
				},
			},
			expectedMoves: []string{"r", "br"},
		},
		
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.clearBoard()
			cfg.IsWhiteTurn = tt.isWhiteTurn
			for _, piece := range tt.pieces {
				cfg.Board[piece.row][piece.col] = Piece{
					Color: piece.color,
					IsKing: piece.isKing,
				}
			}
			moves := cfg.checkSurroundingSquaresForCapture(tt.row, tt.col)
			if !reflect.DeepEqual(moves, tt.expectedMoves) {
				t.Errorf("test checkSurroundingSquaresForCapture() failed. Expected moves: %v   Actual moves: %v", tt.expectedMoves, moves)
			}
		})
	} 
}

func TestKing(t *testing.T) {
	cfg := startCheckers()
	cfg.clearBoard()
	//initialize white piece one away from becoming king
	cfg.Board[1][1] = Piece{
		Color: pieceWhite,
		IsKing: false,
	}

	cfg.movePiece(Move{
		Row: 1,
		Col: 1,
		Direction: moveLeft,
	})

	piece := cfg.Board[0][0]
	if !piece.IsKing {
		t.Errorf("King not working. Piece at (0,0) - Color: %v   IsKing: %v", piece.Color, piece.IsKing)
	}
}

/* func TestDoubleCapture(t *testing.T) {
	cfg := startCheckers()
	cfg.clearBoard()
	cfg.Board[4][4] = Piece{
		Color: pieceWhite,
	}
	cfg.Board[3][3] = Piece{
		Color: pieceBlack,
	}
	cfg.Board[1][1] = Piece{
		Color: pieceBlack,
	}

	cfg.movePiece(Move{
		Row: 4,
		Col: 4,
		Direction: moveLeft,
	})

	if cfg.Board[0][0].Color != pieceWhite {
		t.Error("Double capture failed")
	}
} */