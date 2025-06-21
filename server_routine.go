package main

func StartServerRoutine() {
	cfg := startCheckers()
	serverToClient<- serverToClientData{
		Board: cfg.Board,
		Pieces: cfg.Pieces,
		Error: nil,
		GameOver: false,
	}

	for {
		data := <- clientToServer
		err := cfg.movePiece(data.Move)
		gameOver := false
		if err == nil {
			gameOver = cfg.endTurn()
		}

		serverToClient<- serverToClientData{
			Board: cfg.Board,
			Pieces: cfg.Pieces,
			Error: err,
			GameOver: gameOver,
		}
	}
}