package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
)

const serverURL = "localhost:2000"

func main() {
	fmt.Println("Welcome to checkers!")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\nEnter 1, 2, or 3: ")
	for scanner.Scan() {
		input := scanner.Text()
		switch input{
		case "1":
			StartOnlineGame()
		case "2":
			StartLocalGame()
		case "3":
			GameType = gameSingle
		default:
			fmt.Print("Invalid input.")
		}
		fmt.Print("Please select a game type:\n1. Online Multiplayer\n2. Local Multiplayer\n3. Local Singleplayer\nEnter 1, 2, or 3: ")
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
	serverConn, err := net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Println("Failed to connect to host. Returning to main menu...")
		return
	}
	defer serverConn.Close()
}