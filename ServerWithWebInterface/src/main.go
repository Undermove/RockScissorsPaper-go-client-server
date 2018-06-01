package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Type     string `json:"type"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func provideScriptFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/app.js")
}

func provideStyleFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/style.css")
}

func provideMainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/index.html")
}

func provideRoomPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../room/index.html")
}

func main() {
	// Create a simple file server
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("../public"))
	fs2 := http.FileServer(http.Dir("../room"))
	r.Handle("/", http.StripPrefix("/", fs))
	r.Handle("/room", http.StripPrefix("/room", fs2))
	r.HandleFunc("/app.js", provideScriptFile).Methods("GET")
	r.HandleFunc("/style.css", provideStyleFile).Methods("GET")

	// Configure websocket route
	r.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
