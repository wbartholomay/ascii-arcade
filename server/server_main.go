package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

type Game struct {
	players       [2]net.Conn
	id            int
	playerOneTurn bool
}

var waiting net.Conn

func main() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v", err)
			continue
		}
		go handleNewConnection(conn)
	}
}

func handleNewConnection(conn net.Conn) {
	if waiting == nil {
		waiting = conn
		conn.Write([]byte("1"))
		fmt.Println("Player 1 connected, waiting for Player 2...")
	} else {
		//passing by reference for now
		conn.Write([]byte("2"))
		go StartCheckersGame(&Game{
			players:       [2]net.Conn{waiting, conn},
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
	time.Sleep(2 * time.Second)
	player1.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := player1.Write([]byte("game started"))
	if err != nil {
		fmt.Println("error sending data to client, aborting game.")
		//TODO create a function that shuts down the handleInput go routines, as well as notifies both connections the game has been aborted
	}
	player2.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = player2.Write([]byte("game started"))
	if err != nil {
		fmt.Println("error sending data to client, aborting game.")
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
	
	if err != nil {
		fmt.Println(err)
	}

	for {
		err = currentConn.SendData(checkers.ServerToClientData{
			Board:    cfg.Board,
			Pieces:   cfg.Pieces,
			Error:    nil,
			GameOver: false,
		}, 10)
		if err != nil {
			fmt.Println(err)
		}
		data, err := currentConn.ReceiveData(0)
		if err != nil {
			fmt.Println(err)
		}
		if err = cfg.MovePiece(data.Move, currentConn); err != nil{
			fmt.Println(err)
		}

		if currentConn == &connPlayerOne{
			currentConn = &connPlayerTwo
		} else {
			currentConn = &connPlayerOne
		}
	}

}