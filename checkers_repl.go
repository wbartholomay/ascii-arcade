package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// )

// func StartCheckersRepl(conn net.Conn, playerNum int) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Welcome to ASCII Checkers!")

// 	cfg := startCheckers()
// 	cfg.displayBoard()
// 	fmt.Println("White's Turn:")

// 	for {
// 		fmt.Print("Checkers > ")
// 		scanner.Scan()

// 		t := scanner.Text()
// 		input := cleanInput(t)
// 		if len(input) == 0 { 
// 			continue 
// 		}

// 		cmd, ok := getCommands()[input[0]]
// 		if !ok{
// 			fmt.Println("Unknown command. Enter 'help' to see a list of commands.")
// 			continue
// 		}

// 		//exit checkers game
// 		if cmd.name == "exit" {
// 			break
// 		}

// 		err := cmd.callback(&cfg, input[1:]...)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 	}
// }

// func WaitForServerResponse(conn net.Conn) []byte{
// 	buf := make([]byte, 1024)
// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return buf[:n]
// }