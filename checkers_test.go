package main

import (
	"testing"
)

func TestMovePiece(t *testing.T) {
	tests := []struct {
		name string
		row int8
		col int8
		expectedRow int8
		expectedCol int8
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
			err := movePiece(&cfg, tt.row, tt.col, tt.direction)
			if (err != nil) != tt.wantErr {
				t.Errorf("movePiece() error = %v, wantErr %v. Starting square = (%v, %v), move direction = %v", err, tt.wantErr, tt.row, tt.col, tt.direction)
				return
			}

			if !tt.wantErr && cfg.Board[tt.expectedRow][tt.expectedCol].pieceStatus == pieceEmpty{
				t.Errorf("movePiece() did not move piece to expected square: (%v, %v)", tt.expectedRow, tt.expectedCol)
				return
			}
		})
	}
}