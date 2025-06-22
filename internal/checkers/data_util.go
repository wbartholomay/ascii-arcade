package checkers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Transport[SendT any, RcvT any] interface {
	SendData(SendT, time.Duration) error
	ReceiveData(time.Duration) (RcvT, error)
}

type WebTransport[SendT any, RcvT any] struct {
	Conn net.Conn
}

func (w *WebTransport[SendT, RcvT]) SendData(T SendT, dur time.Duration) error {
	rawData, err := json.Marshal(T)
		if err != nil {
			log.Fatal(err)
		}
		if dur != 0{
			w.Conn.SetWriteDeadline(time.Now().Add(dur * time.Second))
		}
		n, err := w.Conn.Write(rawData)
		if err != nil {
			return err
		}
		if n < len(rawData) {
			return errors.New("not all data was sent to server")
		}
	return nil
}

func (w *WebTransport[SendT, RcvT]) ReceiveData(dur time.Duration) (RcvT, error) {
	//Create 1KB buffer to read from server(could definitely make this smaller, but should not matter)
	// buf := make([]byte, 1024)
	if dur != 0 {
		w.Conn.SetReadDeadline(time.Now().Add(dur * time.Second))
	}
	decoder := json.NewDecoder(w.Conn)
	var data RcvT
	err := decoder.Decode(&data)
	// n, err := w.Conn.Read(buf)
	if err != nil {
		var x RcvT
		return x, fmt.Errorf("error reading data from connection: %w", err)
	}
	// rawData := buf[:n]
	// var data RcvT
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	var x RcvT
	// 	return x, fmt.Errorf("error unmarshaling the json: %w", err)
	// }
	return data, err
}

type LocalTransport[SendT any, RcvT any] struct {
	SendChannel chan SendT
	RcvChannel chan RcvT
}

func (l *LocalTransport[SendT, RcvT]) SendData(T SendT, dur time.Duration) error {
	l.SendChannel <- T
	return nil
}

func (l *LocalTransport[SendT, RcvT]) ReceiveData(dur time.Duration) (RcvT, error) {
	return <- l.RcvChannel, nil
}

type ServerToClientData struct {
	Board             [8][8]Piece    `json:"board"`
	Pieces            map[int]Coords `json:"pieces"`
	IsDoubleJump      bool           `json:"is_double_jump"`
	DoubleJumpOptions []string       `json:"double_jump_options"`
	PieceCoords       [2]int         `json:"piece_coords"`
	Error             error          `json:"error"`
	GameOver          bool           `json:"game_over"`
	NotPlayerTurn     bool           `json:"not_player_turn"`
}

type ClientToServerData struct {
	Move                Move   `json:"move"`
	DoubleJumpDirection string `json:"double_jump_direction"`
}