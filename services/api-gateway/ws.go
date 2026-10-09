package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) error{
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Printf("WebSocket upgrade failed: %v", err)
		return nil
	}
	
	defer conn.Close()
	
	userID := r.URL.Query().Get("userID")
	if userID == ""{	
		log.Println("No user ID provided")
		return nil
	}
	
	for {
		_, message, err := conn.ReadMessage()
		if err != nil{
			log.Panicf("Error reading message: %v",err)
			break
		}
		
		log.Printf("Recived message: %s", message)
	}
	
	return nil
	
}


func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) error{
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Printf("WebSocket upgrade failed: %v", err)
		return nil
	}
	
	defer conn.Close()
	
	userID := r.URL.Query().Get("userID")
	if userID == ""{	
		log.Println("No user ID provided")
		return nil
	}
	
	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == ""{	
		log.Println("No user packageSlug provided")
		return nil
	}
	
	type Driver struct {
		Id               string   `json:"id"`
		Name             string   `json:"name"`
		ProfilePicture   string   `json:"profilePicture"`
		CarPlate         string   `json:"carPlate"`
	    PackgeSlug       string   `json:"packageSlug"`
		
	}
	
	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id: userID,
			Name: "Abdisa Merga",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate: "ABC123",
			PackgeSlug: packageSlug,
		},
	}
	
	if err := conn.WriteJSON(msg); err != nil{
		log.Printf("Error sending message: %v", err)
		return nil
	}
	
	for {
		_, message, err := conn.ReadMessage()
		if err != nil{
			log.Panicf("Error reading message: %v",err)
			break
		}
		
		log.Printf("Recived message: %s", message)
	}
	
	
	return nil
	
}