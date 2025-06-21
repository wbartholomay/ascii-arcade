package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Game struct {
	players [2]net.Conn
	id int
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
		go startCheckersGame(&Game{
			players: [2]net.Conn{waiting, conn},
			id: 0,
			playerOneTurn: true,
		})
		waiting = nil
	}
}


func startCheckersGame(g *Game) {
	fmt.Println("Game started.")
	player1, player2 := g.players[0], g.players[1]
	defer player1.Close()
	defer player2.Close()

	input1 := make(chan string)
	input2 := make(chan string)

	go handleInput(player1, input1, g, true)
	go handleInput(player2, input2, g, false)

	for {
		inputChan := input1
		receivingConn := player2
		if !g.playerOneTurn {
			inputChan = input2
			receivingConn = player1
		}
		input := <- inputChan
		receivingConn.Write([]byte(input))
		g.playerOneTurn = !g.playerOneTurn
	}


}

func handleInput(conn net.Conn, ch chan<- string, g *Game, isPlayerOne bool) {
	scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
		if isPlayerOne != g.playerOneTurn {
			conn.Write([]byte("Waiting on other player...\n"))
			continue
		}
        ch <- scanner.Text()
    }
}