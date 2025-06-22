package main

import (
	"fmt"

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
		//inner for loop to continue requesting moves from the client until no more double jumps are available
		for {
			data, _ := transport.ReceiveData(0)
			nextMoves, pieceCoords, err := cfg.MovePiece(data.Move, &transport)
			hasDoubleJump := len(nextMoves) > 0
			gameOver := false
			if err == nil && !hasDoubleJump{
				gameOver = cfg.EndTurn()
			}

			err = transport.SendData(checkers.ServerToClientData{
				Board:    cfg.Board,
				Pieces:   cfg.Pieces,
				Error:    err,
				GameOver: gameOver,
				IsDoubleJump: hasDoubleJump,
				DoubleJumpOptions: nextMoves,
				PieceCoords: pieceCoords,
			}, 0)
			if err != nil {
				fmt.Println(err)
			}
			if len(nextMoves) == 0 {
				break
			}
			}
	}
}
