package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

func SendDataToServer(data clientToServerData) error{
	if GameType == gameOnline {
		rawData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		serverConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		n, err := serverConn.Write(rawData)
		if err != nil {
			return err
		}
		if n < len(rawData) {
			return errors.New("not all data was sent to server")
		}
	} else {
		clientToServer <- data
	}
	return nil
}

func WaitForDataFromServer() (serverToClientData, error) {
	if GameType == gameOnline {
		//Create 1KB buffer to read from server(could definitely make this smaller, but should not matter)
		buf := make([]byte, 1024)
		//timeout after 10 seconds
		serverConn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := serverConn.Read(buf)
		if err != nil {
			return serverToClientData{}, err
		}
		rawData := buf[:n]
		data := serverToClientData{}
		if err = json.Unmarshal(rawData, &data); err != nil {
			return serverToClientData{}, err
		}
		return data, err
	} else {
		return <- serverToClient, nil
	}
}