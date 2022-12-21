package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var boardcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnection)
	go handleMessage()
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		log.Println("ini MESSAGE", msg)
		boardcast <- msg
	}
}

func handleMessage() {
	for {
		msg := <-boardcast

		for client := range clients {
			// log.Println(client)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error %v", err)
				delete(clients, client)
			}
		}
	}
}
