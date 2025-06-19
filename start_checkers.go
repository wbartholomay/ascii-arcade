package main

func startCheckers() checkersCfg {
	//initialize board
	board := initializeBoard()

	return checkersCfg{
		Board: board,
		IsWhiteTurn: true,
		WhitePieceCount: 12,
		BlackPieceCount: 12,
	}
}

func initializeBoard() [][]Piece {
	board := make([][]Piece, 8)

	whitePieceID := 0
	blackPieceID := 0
	
	for row := range board {
		board[row] = make([]Piece, 8)
		for col := range board[row] {
			hasPiece := ((row % 2) == 0) == ((col % 2) == 0)

			//initialize pieces
			if hasPiece && row < 3 {
				board[row][col] = Piece{
					ID: blackPieceID,
					Color: pieceBlack,
					IsKing: false,
				}
				blackPieceID++
			} else if hasPiece && row > 4 {
				board[row][col] = Piece{
					ID: whitePieceID,
					Color: pieceWhite,
					IsKing: false,
				}
				whitePieceID++
			} else {
				//leaving an empty struct here for now
				board[row][col] = Piece{}
			}
		}
	}

	return board
}