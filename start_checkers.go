package main

func startCheckers() checkersCfg {
	//initialize board
	board, whitePieces, blackPieces := initializeBoard()

	return checkersCfg{
		Board: board,
		WhitePieces: whitePieces,
		BlackPieces: blackPieces,
		IsWhiteTurn: true,
		WhitePieceCount: 12,
		BlackPieceCount: 12,
	}
}

func initializeBoard() ([][]Piece, map[int]Coords, map[int]Coords) {
	board := make([][]Piece, 8)
	whitePieces := map[int]Coords{}
	blackPieces := map[int]Coords{}

	whitePieceID := 1
	blackPieceID := 1
	
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
				blackPieces[blackPieceID] = Coords{
					Row: row,
					Col: col,
				}
				blackPieceID++
			} else if hasPiece && row > 4 {
				board[row][col] = Piece{
					ID: whitePieceID,
					Color: pieceWhite,
					IsKing: false,
				}
				whitePieces[whitePieceID] = Coords{
					Row: row,
					Col: col,
				}
				whitePieceID++
			} else {
				//leaving an empty struct here for now
				board[row][col] = Piece{}
			}
		}
	}

	return board, whitePieces, blackPieces
}