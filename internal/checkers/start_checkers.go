package checkers

func StartCheckers() Checkerscfg {
	//initialize board
	board, pieces := initializeBoard()

	return Checkerscfg{
		Board:           board,
		Pieces:          pieces,
		IsWhiteTurn:     true,
		WhitePieceCount: 12,
		BlackPieceCount: 12,
	}
}

func initializeBoard() ([8][8]Piece, map[int]Coords) {
	board := [8][8]Piece{}
	pieces := map[int]Coords{}

	//start white ids at 101, black ids at 201
	whitePieceID := 101
	blackPieceID := 201

	for row := range board {
		for col := range board[row] {
			hasPiece := ((row % 2) == 0) == ((col % 2) == 0)

			//initialize pieces
			if hasPiece && row < 3 {
				board[row][col] = Piece{
					ID:     blackPieceID,
					Color:  pieceBlack,
					IsKing: false,
				}
				pieces[blackPieceID] = Coords{
					Row: row,
					Col: col,
				}
				blackPieceID++
			} else if hasPiece && row > 4 {
				board[row][col] = Piece{
					ID:     whitePieceID,
					Color:  pieceWhite,
					IsKing: false,
				}
				pieces[whitePieceID] = Coords{
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

	return board, pieces
}
