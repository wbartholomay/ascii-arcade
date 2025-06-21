package main

type serverToClientData struct {
	Board [8][8]Piece `json:"board"`
	Pieces map[int]Coords `json:"pieces"`
	Error error		  `json:"error"`
	GameOver bool     `json:"game_over"`
}

type clientToServerData struct {
	Move Move         `json:"move"`
}

var (
	clientToServer = make(chan clientToServerData)
	serverToClient = make(chan serverToClientData)
)

func main() {
	// fmt.Println("Connecting to server...")
	// conn, err := net.Dial("tcp", "localhost:2000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	// // StartCheckersRepl()

	// buf := make([]byte, 1)
	// _, err = conn.Read(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	go StartServerRoutine()
	ClientRoutine()
}