package checkers

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Transport[SendT any, RcvT any] interface {
	SendData(SendT, time.Duration) error
	ReceiveData(time.Duration) (RcvT, error)
}

type WebTransport[SendT any, RcvT any] struct {
	Conn *websocket.Conn
}

func (w *WebTransport[SendT, RcvT]) SendData(T SendT, dur time.Duration) error {
	if dur != 0{
		w.Conn.SetWriteDeadline(time.Now().Add(dur * time.Second))
	} else {
		w.Conn.SetWriteDeadline(time.Time{})
	}
	return w.Conn.WriteJSON(T)
}

func (w *WebTransport[SendT, RcvT]) ReceiveData(dur time.Duration) (RcvT, error) {
	if dur != 0 {
		w.Conn.SetReadDeadline(time.Now().Add(dur * time.Second))
	} else {
		w.Conn.SetReadDeadline(time.Time{})
	}
	var data RcvT
    err := w.Conn.ReadJSON(&data)
    if err != nil {
        var zero RcvT
        return zero, fmt.Errorf("error reading data from connection: %w", err)
    }
    return data, nil
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
	Error             string          `json:"error"`
	Winner         string           `json:"winner"`
	NotPlayerTurn     bool           `json:"not_player_turn"`
}

type ClientToServerData struct {
	Move                Move   `json:"move"`
	IsConceding				bool   `josn:"is_conceding"`
	DoubleJumpDirection string `json:"double_jump_direction"`
}