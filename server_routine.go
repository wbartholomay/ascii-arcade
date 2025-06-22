package main

import (
	"fmt"

	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

func StartServerRoutine() {
	cfg := checkers.StartCheckers()
	transport := checkers.LocalTransport[checkers.ServerToClientData, checkers.ClientToServerData]{
		SendChannel: serverToClient,
		RcvChannel:  clientToServer,
	}
	transport.SendData(checkers.ServerToClientData{
		Board:  cfg.Board,
		Pieces: cfg.Pieces,
	}, 10)
	for {
		data, err := transport.ReceiveData(0)
		if err != nil {
			fmt.Println(err)
		}
		nextMoves, pieceCoords, moveErr := cfg.MovePiece(data.Move, &transport)
		hasDoubleJump := len(nextMoves) > 0
		errMsg := ""
		if moveErr != nil {
			fmt.Println(moveErr)
			errMsg = moveErr.Error()
		}

		//TODO: gameover is definitely not being passed correctly to clients. SHould add some function which closes connections
		//And invoke that on game over.
		//notify client that of double jump/game over/error stuff
		err = transport.SendData(checkers.ServerToClientData{
			Board:             cfg.Board,
			Pieces:            cfg.Pieces,
			Error:             errMsg,
			Winner:            "",
			IsDoubleJump:      hasDoubleJump,
			DoubleJumpOptions: nextMoves,
			PieceCoords:       pieceCoords,
		}, 5)
		if err != nil {
			fmt.Println(err)
		}

		if !hasDoubleJump && moveErr == nil {
			gameOver := cfg.EndTurn()
			if gameOver {
				EndGame(transport, cfg, cfg.IsWhiteTurn)
			} else {
				err = transport.SendData(checkers.ServerToClientData{
					Board:  cfg.Board,
					Pieces: cfg.Pieces,
					Error:  errMsg,
				}, 5)
				if err != nil {
					fmt.Println(err)
				}
			}

		}

	}
}

func EndGame(t checkers.LocalTransport[checkers.ServerToClientData, checkers.ClientToServerData], cfg checkers.Checkerscfg, whiteWon bool) {
	winner := "w"
	if !whiteWon {
		winner = "b"
	}

	t.SendData(checkers.ServerToClientData{
		Board:  cfg.Board,
		Pieces: cfg.Pieces,
		Winner: winner,
	}, 0)
}
