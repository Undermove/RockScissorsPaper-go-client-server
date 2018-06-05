package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./auth"
	"./data"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

var authModule *auth.AuthModule             // handle clients connections
var broadcast = make(chan WebSocketMessage) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Message msg.Message
	fromWs  *websocket.Conn
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
	authModule = auth.NewModule()

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
	authModule.AddConnection(ws)

	for {
		var msg msg.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			authModule.Disconnect(ws)
			break
		}
		// Send the newly received message to the broadcast channel
		wsMsg := WebSocketMessage{fromWs: ws, Message: msg}
		broadcast <- wsMsg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		wsMsg := <-broadcast

		if authModule.IsLoggedIn(wsMsg.fromWs) {

			// здесь будут расположены какие-то действия в ответ на сообщения, приходящие из сокета
			err := wsMsg.fromWs.WriteJSON(wsMsg.Message)

			if err != nil {
				log.Printf("error: %v", err)
				wsMsg.fromWs.Close()
				authModule.Disconnect(wsMsg.fromWs)
			}
		} else {
			if wsMsg.Message.Type == "AuthRequest" {
				processAuthRequest(wsMsg)
			}
		}
	}
}

func processAuthRequest(wsMsg WebSocketMessage) {
	var request msg.AuthRequest
	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}
	if authModule.Authenticate(wsMsg.fromWs, request) {
		authModule.SendSuccessResponse(wsMsg.fromWs)
	} else {
		authModule.SendRejectResponse(wsMsg.fromWs)
	}
}
