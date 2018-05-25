package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	"./structures"
	"github.com/firstrow/tcp_server"
)

var winMap = map[string]string{
	"Rock":     "Scissors",
	"Paper":    "Rock",
	"Scissors": "Paper",
}

var authConnections = map[string]*tcp_server.Client{}

func main() {
	unknownCounter := 0
	var connections = map[net.Addr]*string{}
	var players = map[string]*structures.Player{}
	var rooms = map[string]*structures.Room{}

	server := tcp_server.New("localhost:9999")

	server.OnNewClient(func(c *tcp_server.Client) {

		var connectionAddress = c.Conn().RemoteAddr()

		if _, ok := connections[connectionAddress]; ok {
			delete(connections, connectionAddress)
		}

		unknownCounter++
		unknownUser := "unknown#" + strconv.Itoa(unknownCounter)
		connections[connectionAddress] = &unknownUser
	})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		var detachedMessage = strings.Split(message, ";")

		currentUserName := detachedMessage[0]
		command := detachedMessage[1]

		if currentUserName == "unknownUser" {
			if command == "reg" {
				if registerUser(detachedMessage[2], c.Conn().RemoteAddr(), players, connections) == false {
					c.Send("User with such name already registred.")
					return
				}
				authConnections[detachedMessage[2]] = c
				c.Send("registred")
				c.Send("Congradulations! You have been registred. Your next command:\n1. rooms - watch all rooms\n2. newroom;name - create new room with name\n3. enter;name - create new room with name")
			} else {
				c.Send("Can't listen commands from unknown user")
			}
		} else if _, ok := players[currentUserName]; ok {
			if command == "rooms" {
				getAllRoomsNames(rooms, c)
			} else if command == "newroom" {
				registerRoom(rooms, detachedMessage[2], players[currentUserName])
			} else if command == "enter" {
				isEntered := enterRoom(rooms, players[currentUserName].Name, detachedMessage[2], players)
				if isEntered == true {
					c.Send("Game started\nAvaliable commands: \n1. turn;yourChoise - example turn;scissors\n2. leave - Leave room")
				} else {
					c.Send("Room not found")
				}
			}

			if rooms[players[currentUserName].CurrentRoomName] != nil {
				if command == "leave" {
					leaveRoom(rooms, currentUserName, players)
				} else if command == "turn" {
					turn(rooms, currentUserName, players, detachedMessage[2], connections)
				}
			}
		}
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		connectionAddress := c.Conn().RemoteAddr()
		username := *connections[connectionAddress]
		delete(players, username)
		delete(connections, connectionAddress)
		log.Println("Disconnected")
	})

	server.Listen()
}

func registerUser(username string,
	address net.Addr,
	players map[string]*structures.Player,
	connections map[net.Addr]*string) bool {
	if _, ok := players[username]; ok {
		return false
	}

	players[username] = structures.NewPlayer(username)
	connections[address] = &username
	return true
}

func getAllRoomsNames(rooms map[string]*structures.Room, c *tcp_server.Client) {
	for _, value := range rooms {
		playerinfo := " |players: "
		for i := 0; i < 2; i++ {
			if value.Players[i] != nil {
				playerinfo = playerinfo + value.Players[i].Name + " "
			}
		}
		c.Send(value.Name + playerinfo)
	}
}

func registerRoom(rooms map[string]*structures.Room, name string, player *structures.Player) {
	rooms[name] = structures.NewRoom(name, player)
}

func enterRoom(rooms map[string]*structures.Room, playername string, roomname string, players map[string]*structures.Player) bool {
	if _, ok := rooms[roomname]; ok {
		rooms[roomname].Players[1] = players[playername]
		players[playername].CurrentRoomName = roomname
		return true
	}

	return false
}

func leaveRoom(rooms map[string]*structures.Room, playername string, players map[string]*structures.Player) bool {
	for i := 0; i < 2; i++ {

		if rooms[players[playername].CurrentRoomName].Players[i].Name == playername {
			rooms[players[playername].CurrentRoomName].Players[i] = structures.NewPlayer("")
			return true
		}
	}

	return false
}

func turn(rooms map[string]*structures.Room,
	playername string,
	players map[string]*structures.Player,
	playerChoise string,
	connections map[net.Addr]*string) {
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
	} else if winMap[room.Players[0].PlayerChoise] == room.Players[1].PlayerChoise {
		result = room.Players[0].Name + " WINS!!!"
	} else {
		result = room.Players[1].Name + " WINS!!!"
	}

	for _, currentPlayer := range room.Players {
		authConnections[currentPlayer.Name].Send(result)
	}
}
