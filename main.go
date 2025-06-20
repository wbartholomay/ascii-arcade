package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type gameType int

const (
	gameLocal gameType = iota
	gameOnline
	gameSingle
)

var GameType gameType

type serverToClientData struct {
	Board [8][8]Piece `json:"board"`
	Pieces map[int]Coords `json:"pieces"`
	IsDoubleJump bool  `json:"is_double_jump"`
	DoubleJumpOptions []string `json:"double_jump_options"`
	Error error		  `json:"error"`
	GameOver bool     `json:"game_over"`
}

type clientToServerData struct {
	Move Move         `json:"move"`
	DoubleJumpDirection string `json:"double_jump_direction"`
}

var (
	clientToServer chan clientToServerData
	serverToClient chan serverToClientData
	serverConn net.Conn
	playerNumber string
)

const serverURL = "localhost:2000"

func main() {
	fmt.Println("Welcome to checkers!")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\n4. Exit\nEnter 1, 2, 3, or 4: ")
	for scanner.Scan() {
		input := scanner.Text()
		switch input{
		case "1":
			StartOnlineGame()
		case "2":
			StartLocalGame()
		case "3":
			GameType = gameSingle
		case "4":
			fmt.Println("See ya later!")
			os.Exit(0)
		default:
			fmt.Print("Invalid input. ")
		}
		fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\n4. Exit\nEnter 1, 2, 3, or 4: ")
	}
	// buf := make([]byte, 1)
	// _, err = conn.Read(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
}

func StartLocalGame() {
	GameType = gameLocal
	serverToClient = make(chan serverToClientData)
	clientToServer = make(chan clientToServerData)
	go StartServerRoutine()
	ClientRoutine()
}

func StartOnlineGame() {
	GameType = gameOnline
	fmt.Println("Connecting to server...")

	//this local serverConn is overwriting global server conn when using := operator
	var err error
	serverConn, err = net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Println("Failed to connect to host. Returning to main menu...")
		return
	}
	//TODO: create abstraction of this function, which sends something to the server to notify the client is closed?
	//May not be necessary as the server will timeout when trying to read from the client, will see about this
	defer serverConn.Close()

	buf := make([]byte, 1024)
	serverConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = serverConn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	playerNumber = string(buf[0])
	fmt.Printf("You are player %v, waiting to start game...\n",playerNumber)
	//TODO: add ability for player to exit game while waiting for game to start (without closing program)
	serverConn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	n, err := serverConn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(buf[:n]) != "game started" {
		fmt.Println("Received unexpected message from server, closing game.")
		return
	}

	fmt.Println("Opponent found, starting game!")
	ClientRoutine()
}