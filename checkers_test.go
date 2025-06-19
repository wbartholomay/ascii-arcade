package main

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			cfg := startCheckers()
			err := cfg.movePiece(Move{
				Row: tt.row,
				Col: tt.col,
				Direction: tt.direction,
				DestRow: tt.row,
				DestCol: tt.col,
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
				DestRow: tt.row,
				DestCol: tt.col,
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("movePiece() error = %v, wantErr %v. Starting square = (%v, %v), move direction = %v", err, tt.wantErr, tt.row, tt.col, tt.direction)
			}

			if !tt.wantErr && cfg.Board[tt.expectedRow][tt.expectedCol].Color == ""{
				t.Errorf("movePiece() did not move piece to expected square: (%v, %v)", tt.expectedRow, tt.expectedCol)
				return
			}
		})
	}
}

func TestApplyDirection(t *testing.T) {
	tests := []struct {
		name string
		move Move
		destRow int
		destCol int
	} {
		{
			name: "Test moving left",
			move: Move{
				Row: 1,
				Col: 1,
				Direction: moveLeft,
				DestRow: 1,
				DestCol: 1,
			},
			destRow: 0,
			destCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			tt.move.applyDirection()
			if tt.move.DestRow != tt.destRow || tt.move.DestCol != tt.destCol {
				t.Errorf("applyDirection() failed - expected coords (%v, %v) - actual coords (%v, %v)", tt.move.DestRow, tt.move.DestCol, tt.destRow, tt.destCol)
			}
		})
	}
}