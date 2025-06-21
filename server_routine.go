package main

import (
	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

func StartServerRoutine() {
	cfg := checkers.StartCheckers()
	serverToClient <- checkers.ServerToClientData{
		Board:    cfg.Board,
		Pieces:   cfg.Pieces,
		Error:    nil,
		GameOver: false,
	}

	for {
		data := <-clientToServer
		err := cfg.MovePiece(data.Move, serverToClient, clientToServer)
		gameOver := false
		if err == nil {
			gameOver = cfg.EndTurn()
		}

		serverToClient <- checkers.ServerToClientData{
			Board:    cfg.Board,
			Pieces:   cfg.Pieces,
			Error:    err,
			GameOver: gameOver,
		}
	}
}
