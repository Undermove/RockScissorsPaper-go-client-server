package main

import (
	"log"
	"strconv"
	"strings"

	"./structures"
	"github.com/firstrow/tcp_server"
)

var loosersMap = map[string]string{
	"Rock":     "Scissors",
	"Paper":    "Rock",
	"Scissors": "Paper",
}

var connections = map[*tcp_server.Client]string{}
var authConnections = map[string]*tcp_server.Client{}
var rooms = map[string]*structures.Room{}
var players = map[string]*structures.Player{}

func main() {
	unknownCounter := 0

	server := tcp_server.New("localhost:9999")

	server.OnNewClient(func(c *tcp_server.Client) {
		unknownCounter++
		unknownUser := "unknown#" + strconv.Itoa(unknownCounter)
		connections[c] = unknownUser
	})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		var detachedMessage = strings.Split(message, ";")

		userName := detachedMessage[0]
		command := detachedMessage[1]
		argument := detachedMessage[2]

		if player, ok := tryGetUser(userName); ok {
			if command == "rooms" {
				getAllRoomsNames(c)
			} else if command == "newroom" {
				registerRoom(argument, player)
			} else if command == "enter" {
				enterRoom(player.Name, argument)
			}

			if rooms[player.CurrentRoomName] != nil {
				if command == "leave" {
					leaveRoom(userName)
				} else if command == "turn" {
					turn(userName, argument)
				}
			}
		} else {
			if command == "reg" {
				registerUser(argument, c)
			} else {
				c.Send("Can't listen commands from unknown user. Register first")
			}
		}
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		username := connections[c]
		leaveRoom(username)
		delete(players, username)
		delete(authConnections, username)
		delete(connections, c)

		log.Println("Disconnected")
	})

	server.Listen()
}

func registerUser(username string, c *tcp_server.Client) {
	if _, ok := tryGetUser(username); ok {
		c.Send("User with such name already registred.")
		return
	}

	players[username] = structures.NewPlayer(username)
	authConnections[username] = c
	connections[c] = username
	c.Send("registred")
	c.Send("Congradulations! You have been registred. Your next command:\n   rooms - watch all rooms\n   newroom;name - create new room with name\n   enter;name - create new room with name")
	return
}

func getAllRoomsNames(c *tcp_server.Client) {
	result := "Rooms:\n"
	for _, value := range rooms {
		playerinfo := " |players: "
		for i := 0; i < 2; i++ {
			if value.Players[i] != nil {
				playerinfo = playerinfo + value.Players[i].Name + " "
			}
		}
		result = result + value.Name + playerinfo
	}

	c.Send(result)
}

func registerRoom(name string, player *structures.Player) {
	rooms[name] = structures.NewRoom(name, player)
	authConnections[player.Name].Send("Room registred")
}

func enterRoom(playername string, roomname string) {
	conn := authConnections[playername]

	if _, ok := rooms[roomname]; ok {
		if players[playername].CurrentRoomName == roomname {
			conn.Send("You are already in this room")
		}

		if len(rooms[roomname].Players) >= 2 {
			conn.Send("Room is full")
		}

		rooms[roomname].Players[1] = players[playername]
		players[playername].CurrentRoomName = roomname
		conn.Send("Game started\nAvaliable commands: \n   turn;yourChoise - example turn;scissors\n   leave - Leave room")
		return
	}

	conn.Send("Room not found")
}

func leaveRoom(playername string) bool {
	if players[playername].CurrentRoomName == "" {
		return false
	}

	for i := 0; i < 2; i++ {
		currentPlayer := rooms[players[playername].CurrentRoomName].Players[i]
		if currentPlayer.Name == playername {
			rooms[players[playername].CurrentRoomName].Players[i] = structures.NewPlayer("")
			authConnections[playername].Send("Room left")
			return true
		}
	}

	return false
}

func turn(playername string, playerChoise string) {
	player := players[playername]
	player.SetPlayerChoise(playerChoise)

	room := rooms[player.CurrentRoomName]
	for _, currentPlayer := range room.Players {
		if currentPlayer.PlayerChoise == "" {
			return
		}
	}

	var result string

	if room.Players[0].PlayerChoise == room.Players[1].PlayerChoise {
		result = "DRAW"
	} else if loosersMap[room.Players[0].PlayerChoise] == room.Players[1].PlayerChoise {
		result = room.Players[0].Name + " WINS!!!"
	} else {
		result = room.Players[1].Name + " WINS!!!"
	}

	for _, currentPlayer := range room.Players {
		authConnections[currentPlayer.Name].Send(result)
	}
}

func tryGetUser(userName string) (*structures.Player, bool) {
	player, ok := players[userName]
	return player, ok
}
