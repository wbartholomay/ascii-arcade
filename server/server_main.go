package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Game struct {
	players       [2]*websocket.Conn
	id            int
	playerOneTurn bool
}

var waiting *websocket.Conn

func main() {
	http.HandleFunc("/ws", handleWSConnection)

	fmt.Println("Server starting on :2000...")
	if err := http.ListenAndServe(":2000", nil); err != nil {
		log.Fatal(err)
	}
}

func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v", err)
		return
	}
	
	handleNewConnection(conn)
}

func handleNewConnection(conn *websocket.Conn) {
	type playerNum struct {
		PlayerNumber string `json:"player_number"`
	}
	if waiting == nil {
		waiting = conn
		err := conn.WriteJSON(playerNum{
			PlayerNumber: "1",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Player 1 connected, waiting for Player 2...")
	} else {
		err := conn.WriteJSON(playerNum{
			PlayerNumber: "2",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		go StartCheckersGame(&Game{
			players:       [2]*websocket.Conn{waiting, conn},
			id:            0,
			playerOneTurn: true,
		})
		waiting = nil
	}
}

func StartCheckersGame(g *Game) {
	fmt.Println("Game started.")
	player1, player2 := g.players[0], g.players[1]
	defer player1.Close()
	defer player2.Close()

	//signal to clients that the game has started
	time.Sleep(1 * time.Second)
	type gameStart struct {
		GameStart bool `json:"game_start"`
	}
	player1.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err := player1.WriteJSON(gameStart{
		GameStart: true,
	})
	if err != nil {
		fmt.Println("error sending data to client, aborting game.")
		os.Exit(0)
		//TODO create a function that shuts down the handleInput go routines, as well as notifies both connections the game has been aborted
	}
	player2.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = player2.WriteJSON(gameStart{
		GameStart: true,
	})
	if err != nil {
		fmt.Println("error sending data to client, aborting game.")
		os.Exit(0)
		//TODO create a function that shuts down the handleInput go routines, as well as notifies both connections the game has been aborted
	}

	connPlayerOne := checkers.WebTransport[checkers.ServerToClientData, checkers.ClientToServerData] {
		Conn: player1,
	}
	connPlayerTwo := checkers.WebTransport[checkers.ServerToClientData, checkers.ClientToServerData] {
		Conn: player2,
	}
	currentConn := &connPlayerOne

	cfg := checkers.StartCheckers()

	err = currentConn.SendData(checkers.ServerToClientData{
		Board:    cfg.Board,
		Pieces:   cfg.Pieces,
		Error:    "",
		GameOver: false,
	}, 10)
	if err != nil {
		fmt.Println(err)
	}

	for {
			data, err := currentConn.ReceiveData(0)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Client disconnected, shutting down.")
				break
			}
			nextMoves, pieceCoords, moveErr := cfg.MovePiece(data.Move, currentConn)
			hasDoubleJump := len(nextMoves) > 0
			gameOver := false
			errMsg := ""
			if moveErr != nil {
				fmt.Println(moveErr)
				errMsg = moveErr.Error()
			}

			//TODO: gameover is definitely not being passed correctly to clients. SHould add some function which closes connections
			//And invoke that on game over.
			//notify client that of double jump/game over/error stuff
			err = currentConn.SendData(checkers.ServerToClientData{
				Board:    cfg.Board,
				Pieces:   cfg.Pieces,
				Error:    errMsg,
				GameOver: gameOver,
				IsDoubleJump: hasDoubleJump,
				DoubleJumpOptions: nextMoves,
				PieceCoords: pieceCoords,
			}, 5)
			if err != nil {
				fmt.Println(err)
			}

			//connections should not be swapped on double jumps and failed moves
			if !hasDoubleJump && moveErr == nil{
				gameOver = cfg.EndTurn()
				//switch conn
				if currentConn == &connPlayerOne {
					currentConn = &connPlayerTwo
				} else {
					currentConn = &connPlayerOne
				}

				err = currentConn.SendData(checkers.ServerToClientData{
					Board:    cfg.Board,
					Pieces:   cfg.Pieces,
					Error:    errMsg,
					GameOver: gameOver,
				}, 5)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Client disconnected, shutting down.")
					break
				}
			}

	}

}