package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/wbarthol/ascii-arcade/internal/checkers"
)

type clientData struct {
	Pieces             map[int]checkers.Coords
	IsWhiteTurn        bool
}

func ClientRoutine(transport checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData]) {
	//TODO: TIMING OUT HERE. NEED TO CONFIGURE SENDING DATA FROM THE SERVER TO THIS POINT.
	whiteTurn := true
	if playerNumber == "2" {
		whiteTurn = false
		fmt.Println("Waiting for player 1 to make their move...")
	}
	fmt.Println()


	for {
		data, err := transport.ReceiveData(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		cfg := clientData{
			Pieces:             data.Pieces,
			IsWhiteTurn:        whiteTurn,
		}

		checkers.DisplayBoard(data.Board, cfg.IsWhiteTurn)
		
		repl:
		for {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Checkers > ")
			scanner.Scan()

			t := scanner.Text()
			input := cleanInput(t)
			if len(input) == 0 {
				continue
			}

			cmd := input[0]
			switch cmd{
			case "move":
				move, err := validateMove(&cfg, input[1:]...)
				if err != nil {
					fmt.Printf("Error validating move: %v\n", err)
					continue
				}

				err = sendMoveToServer(&cfg, transport, move)
				if err != nil {
					fmt.Printf("Error when sending move to server: %v\n", err)
					continue
				}

				break repl
			case "help":
				commandHelp()
			case "concede":
				commandConcede(cfg.IsWhiteTurn)
				//TODO exit to main menu on concession
				os.Exit(0)
			}
		}
	}
}

type cliCommand struct {
	name        string
	description string
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a list of commmands",
		},
		"move": {
			name:        "move",
			description: "Move a piece. Takes arguments <piece-number> <direction {'l', 'r', 'bl', 'br'}>",
		},
		"concede": {
			name:        "concede",
			description: "Concede the game",
		},
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	substrings := strings.Fields(text)
	return substrings
}

func validateMove(cfg *clientData, 
	params ...string) (checkers.Move, error) {
	//validate params - expecting move <row> <col> <direction>
	if len(params) < 2 {
		return checkers.Move{}, errors.New("not enough arguments. Expecting move <piece number> <direction>")
	}

	pieceNum, err := strconv.ParseInt(params[0], 10, 8)
	if err != nil {
		return checkers.Move{}, fmt.Errorf("expected a number, got: %v", params[0])
	}

	pieceId := checkers.GetActualID(checkers.GetPlayerColor(cfg.IsWhiteTurn), int(pieceNum))
	piece, ok := cfg.Pieces[pieceId]
	if !ok {
		return checkers.Move{}, fmt.Errorf("invalid piece number: %v", params[0])
	}

	directionString := params[1]
	direction, ok := checkers.MovesMap[directionString]
	if !ok {
		return checkers.Move{}, errors.New("invalid move direction. Valid moves are 'l', 'r', 'bl', 'br'")
	}

	return checkers.Move{
		Row:       piece.Row,
		Col:       piece.Col,
		Direction: direction,
	}, nil
}

func sendMoveToServer(clientData *clientData,
	T checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData], 
	move checkers.Move) error {
	//send move to server. TODO replace this and other sending of data with abstractions which check the game type
	err := T.SendData(checkers.ClientToServerData{
		Move: move,
	}, 10)
	if err != nil {
		return err
	}

	data := checkers.ServerToClientData{}
	for {
		data, err = T.ReceiveData(5)
		if err != nil {
			return err
		}
		if data.Error != "" {
			return errors.New(data.Error)
		}
		//DOUBLE JUMP HANDLING: MAYBE COULD MOVE THIS OUT INTO MAIN FUNC?
		if data.IsDoubleJump {
			checkers.DisplayBoard(data.Board, clientData.IsWhiteTurn)
			fmt.Print("Another capture is available, enter one of the following directions: ")
			for _, moveStr := range data.DoubleJumpOptions {
				fmt.Printf("%v, ", moveStr)
			}
			fmt.Println()
			input := ""
			for {
				fmt.Print("Checkers > ")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input = strings.ToLower(scanner.Text())
				if !slices.Contains(data.DoubleJumpOptions, input) {
					fmt.Println("Please enter one of the displayed directions.")
				} else {
					break
				}
			}
			err = T.SendData(checkers.ClientToServerData{
				Move: checkers.Move{
					Row: data.PieceCoords[0],
					Col: data.PieceCoords[1],
					Direction: checkers.MovesMap[input],
				},
			}, 5)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	if GameType != gameOnline {
		clientData.IsWhiteTurn = !clientData.IsWhiteTurn
	}

	clientData.Pieces = data.Pieces

	if data.GameOver {
		return errors.New("game over")
	}

	return nil
}

func commandHelp() {
	fmt.Print("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
}

func commandConcede(isWhiteTurn bool) {
	if isWhiteTurn {
		fmt.Println("White conceded, black wins!")
	} else {
		fmt.Println("Black conceded, white wins!")
	}
}
