package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

type gameType int

const (
	gameLocal gameType = iota
	gameOnline
	gameSingle
)

var GameType gameType

// definitely overusing global scope, should remove these
var (
	clientToServer chan checkers.ClientToServerData
	serverToClient chan checkers.ServerToClientData
	playerNumber   string
)

// TODO: set this as environment variable
const serverURL = "localhost:2000/ws"

func main() {
	fmt.Println("Welcome to checkers!")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\n4. Exit\nEnter 1, 2, 3, or 4: ")
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "1":
			StartOnlineGame()
		case "2":
			StartLocalGame()
		case "3":
			GameType = gameSingle
			fmt.Println("Local singleplayer coming soon!")
		case "4":
			fmt.Println("See ya later!")
			os.Exit(0)
		default:
			fmt.Print("Invalid input. ")
		}
		fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\n4. Exit\nEnter 1, 2, 3, or 4: ")
	}
}

func StartLocalGame() {
	GameType = gameLocal
	serverToClient = make(chan checkers.ServerToClientData)
	clientToServer = make(chan checkers.ClientToServerData)
	go StartServerRoutine()
	transport := checkers.LocalTransport[checkers.ClientToServerData, checkers.ServerToClientData]{
		SendChannel: clientToServer,
		RcvChannel:  serverToClient,
	}
	ClientRoutine(&transport)
}

func StartOnlineGame() {
	GameType = gameOnline
	fmt.Println("Connecting to server...")

	//this local serverConn is overwriting global server conn when using := operator
	var err error
	serverConn, _, err := websocket.DefaultDialer.Dial("ws://"+serverURL, nil)
	if err != nil {
		fmt.Println("Failed to connect to host. Returning to main menu...")
		return
	}
	//TODO: create abstraction of this function, which sends something to the server to notify the client is closed?
	//May not be necessary as the server will timeout when trying to read from the client, will see about this
	defer serverConn.Close()
	transport := checkers.WebTransport[checkers.ClientToServerData, checkers.ServerToClientData]{
		Conn: serverConn,
	}

	type playerNum struct {
		PlayerNumber string `json:"player_number"`
	}
	var buf playerNum
	serverConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	err = serverConn.ReadJSON(&buf)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to connect to host. Returning to main menu...")
		return
	}
	playerNumber = buf.PlayerNumber
	fmt.Printf("You are player %v, waiting to start game...\n", playerNumber)
	//TODO: add ability for player to exit game while waiting for game to start (without closing program)
	type gameStart struct {
		GameStart bool `json:"game_start"`
	}
	var g gameStart
	serverConn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	err = serverConn.ReadJSON(&g)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Received unexpected message from server, closing game.")
		return
	}
	if !g.GameStart {
		fmt.Println("Received unexpected message from server, closing game.")
		return
	}

	fmt.Println("Opponent found, starting game!")
	ClientRoutine(&transport)
}
