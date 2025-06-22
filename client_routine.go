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

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Checkers > ")
			scanner.Scan()

			t := scanner.Text()
			input := cleanInput(t)
			if len(input) == 0 {
				continue
			}

			cmd, ok := getCommands()[input[0]]
			if !ok {
				fmt.Println("Unknown command. Enter 'help' to see a list of commands.")
				continue
			}

			err := cmd.callback(&cfg, transport, input[1:]...)
			//using this error to pass game over state, might need to update that
			if err != nil {
				if err.Error() == "game over" {
					break
				}
				fmt.Println(err.Error())
				continue
			}
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *clientData, 
		transport checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData],
		 params ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a list of commmands",
			callback:    commandHelp,
		},
		"move": {
			name:        "move",
			description: "Move a piece. Takes arguments <piece-number> <direction {'l', 'r', 'bl', 'br'}>",
			callback:    commandMove,
		},
		"concede": {
			name:        "concede",
			description: "Concede the game",
			callback:    commandConcede,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	substrings := strings.Fields(text)
	return substrings
}

func commandMove(cfg *clientData, 
	T checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData], 
	params ...string) error {
	//validate params - expecting move <row> <col> <direction>
	if len(params) < 2 {
		return errors.New("not enough arguments. Expecting move <piece number> <direction>")
	}

	pieceNum, err := strconv.ParseInt(params[0], 10, 8)
	if err != nil {
		return fmt.Errorf("expected a number, got: %v", params[0])
	}

	pieceId := checkers.GetActualID(checkers.GetPlayerColor(cfg.IsWhiteTurn), int(pieceNum))
	piece, ok := cfg.Pieces[pieceId]
	if !ok {
		return fmt.Errorf("invalid piece number: %v", params[0])
	}

	directionString := params[1]
	direction, ok := checkers.MovesMap[directionString]
	if !ok {
		return errors.New("invalid move direction. Valid moves are 'l', 'r', 'bl', 'br'")
	}

	move := checkers.Move{
		Row:       piece.Row,
		Col:       piece.Col,
		Direction: direction,
	}

	//send move to server. TODO replace this and other sending of data with abstractions which check the game type
	err = T.SendData(checkers.ClientToServerData{
		Move: move,
	}, 10)
	if err != nil {
		return err
	}

	data := checkers.ServerToClientData{}
	for {
		data, err = T.ReceiveData(10)
		if err != nil {
			return err
		}
		if data.IsDoubleJump {
			checkers.DisplayBoard(data.Board, cfg.IsWhiteTurn)
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
			}, 10)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	if data.Error != nil {
		return data.Error
	}

	if GameType != gameOnline {
		cfg.IsWhiteTurn = !cfg.IsWhiteTurn
	}

	cfg.Pieces = data.Pieces

	if data.GameOver {
		return errors.New("game over")
	}
	return nil
}

func commandHelp(cfg *clientData,
	T checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData],
	 params ...string) error {
	fmt.Print("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandConcede(cfg *clientData,
	T checkers.Transport[checkers.ClientToServerData, checkers.ServerToClientData],
	params ...string) error {
	if cfg.IsWhiteTurn {
		fmt.Println("White conceded, black wins!")
	} else {
		fmt.Println("Black conceded, white wins!")
	}

	return errors.New("game over")
}
