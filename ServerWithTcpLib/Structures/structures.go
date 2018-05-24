package structures

import (
	"net"
)

// Room with players
type Room struct {
	Name    string
	Players [2]*Player
}

func NewRoom(name string, player *Player) *Room {
	var players [2]*Player
	players[0] = player

	room := &Room{
		Name:    name,
		Players: players,
	}

	return room
}

// Player contains information about user and his connection
type Player struct {
	Name         string
	Address      net.Addr
	Score        int
	PlayerChoise string
}

func (player *Player) SetPlayerChoise(choise string) {
	player.PlayerChoise = choise
}

func NewPlayer(name string) *Player {
	player := &Player{
		Name:         name,
		PlayerChoise: "",
		Score:        0,
	}

	return player
}
