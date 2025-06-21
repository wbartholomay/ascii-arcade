package main

import (
	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

func StartServerRoutine() {
	cfg := checkers.StartCheckers()
	transport := checkers.LocalTransport[checkers.ServerToClientData, checkers.ClientToServerData]{
		SendChannel: serverToClient,
		RcvChannel: clientToServer,
	}
	transport.SendData(checkers.ServerToClientData{
		Board:    cfg.Board,
		Pieces:   cfg.Pieces,
		Error:    nil,
		GameOver: false,
	}, 10)

	for {
		data := <-clientToServer
		err := cfg.MovePiece(data.Move, &transport)
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
