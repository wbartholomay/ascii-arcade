package checkers

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
)

type Transport[SendT any, RcvT any] interface {
	SendData(SendT) error
	ReceiveData() (RcvT, error)
}

type WebTransport[SendT any, RcvT any] struct {
	Conn net.Conn
}

func (w *WebTransport[SendT, RcvT]) SendData(T SendT) error {
	rawData, err := json.Marshal(T)
		if err != nil {
			log.Fatal(err)
		}
		w.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := w.Conn.Write(rawData)
		if err != nil {
			return err
		}
		if n < len(rawData) {
			return errors.New("not all data was sent to server")
		}
	return nil
}

func (w *WebTransport[SendT, RcvT]) ReceiveData() (RcvT, error) {
	//Create 1KB buffer to read from server(could definitely make this smaller, but should not matter)
	buf := make([]byte, 1024)
	//timeout after 10 seconds
	w.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := w.Conn.Read(buf)
	if err != nil {
		var x RcvT
		return x, err
	}
	rawData := buf[:n]
	var data RcvT
	if err = json.Unmarshal(rawData, &data); err != nil {
		var x RcvT
		return x, err
	}
	return data, err
}

type LocalTransport[SendT any, RcvT any] struct {
	SendChannel chan SendT
	RcvChannel chan RcvT
}

func (l *LocalTransport[SendT, RcvT]) SendData(T SendT) error {
	l.SendChannel <- T
	return nil
}

func (l *LocalTransport[SendT, RcvT]) ReceiveData() (RcvT, error) {
	return <- l.RcvChannel, nil
}

type ServerToClientData struct {
	Board             [8][8]Piece    `json:"board"`
	Pieces            map[int]Coords `json:"pieces"`
	IsDoubleJump      bool           `json:"is_double_jump"`
	DoubleJumpOptions []string       `json:"double_jump_options"`
	Error             error          `json:"error"`
	GameOver          bool           `json:"game_over"`
	NotPlayerTurn     bool           `json:"not_player_turn"`
}

type ClientToServerData struct {
	Move                Move   `json:"move"`
	DoubleJumpDirection string `json:"double_jump_direction"`
}

// Send data to server - accepts data and either a network connection or a channel
/* func SendDataToServer(data ClientToServerData, channel interface{}) error {
	switch c := channel.(type) {
	case net.Conn:
		rawData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := c.Write(rawData)
		if err != nil {
			return err
		}
		if n < len(rawData) {
			return errors.New("not all data was sent to server")
		}
	case chan ClientToServerData:
		c <- data
	default:
		return errors.New("server not recognized")
	}
	return nil
}

// Wait for data from server - accepts a network connection or channel
func WaitForDataFromServer(channel interface{}) (ServerToClientData, error) {
	switch c := channel.(type) {
	case net.Conn:
		//Create 1KB buffer to read from server(could definitely make this smaller, but should not matter)
		buf := make([]byte, 1024)
		//timeout after 10 seconds
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := c.Read(buf)
		if err != nil {
			return ServerToClientData{}, err
		}
		rawData := buf[:n]
		data := ServerToClientData{}
		if err = json.Unmarshal(rawData, &data); err != nil {
			return ServerToClientData{}, err
		}
		return data, err
	case chan ServerToClientData:
		return <-c, nil
	default:
		return ServerToClientData{}, errors.New("server not recognized")
	}
}

// TODO: could merge these functions using an interface, not sure what functions the types should implement though
// Send data to server - accepts data and either a network connection or a channel
func SendDataToClient(data ServerToClientData, channel interface{}) error {
	switch c := channel.(type) {
	case net.Conn:
		rawData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := c.Write(rawData)
		if err != nil {
			return err
		}
		if n < len(rawData) {
			return errors.New("not all data was sent to server")
		}
	case chan ServerToClientData:
		c <- data
	default:
		return errors.New("server not recognized")
	}
	return nil
}

// Wait for data from server - accepts a network connection or channel
func WaitForDataFromClient(channel interface{}) (ClientToServerData, error) {
	switch c := channel.(type) {
	case net.Conn:
		//Create 1KB buffer to read from server(could definitely make this smaller, but should not matter)
		buf := make([]byte, 1024)
		//timeout after 10 seconds
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := c.Read(buf)
		if err != nil {
			return ClientToServerData{}, err
		}
		rawData := buf[:n]
		data := ClientToServerData{}
		if err = json.Unmarshal(rawData, &data); err != nil {
			return ClientToServerData{}, err
		}
		return data, err
	case chan ClientToServerData:
		return <-c, nil
	default:
		return ClientToServerData{}, errors.New("server not recognized")
	}
}
*/