package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	"./structures"
	"github.com/firstrow/tcp_server"
)

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

		if detachedMessage[0] == "unknownUser" {
			if detachedMessage[1] == "reg" {
				if registerUser(detachedMessage[2], c.Conn().RemoteAddr(), players, connections) == false {
					c.Send("User with such name already registred.")
					return
				}

				c.Send("registred")
				c.Send("Congradulations! You have been registred. Your next command:\n1. rooms - watch all rooms\n2. newroom;name - create new room with name\n3. enter;name - create new room with name")
			} else {
				c.Send("Can't listen commands from unknown user")
			}
		} else if _, ok := players[detachedMessage[0]]; ok {
			command := detachedMessage[1]
			if command == "rooms" {
				getAllRoomsNames(rooms, c)
			} else if command == "newroom" {
				registerRoom(rooms, detachedMessage[2], players[detachedMessage[0]])
			} else if command == "enter" {

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

func enterRoom(rooms map[string]*structures.Room, name string, players map[string]*structures.Player) {
	rooms[name].Players[1] = players[name]
}
