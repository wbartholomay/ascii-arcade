package main

import (
	"encoding/json"
	"log"
)

func SendDataToServer(data clientToServerData) {
	if GameType == gameOnline {
		_, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		//TODO: ADD LOGIC TO SEND DATA TO SERVER THROUGH CONNECTION (CONNECTION SHOULD GET DEFINED IN GLOBAL SCOPE)

	} else {
		clientToServer <- data
	}
}