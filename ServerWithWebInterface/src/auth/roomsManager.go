package auth

import "github.com/gorilla/websocket"

// Room with players
type Room struct {
	Name    string
	Players [2]*Player
}

func (r *Room) EnterRoom(playerName string) {
	if len(r.Players) <= 2 {

	}
}

func (r Room) LeaveRoom(playerName string) {

}

type Player struct {
	Name   string
	Score  int
	Choise string
}

type RoomsManager struct {
	Rooms       map[string]*Room
	ConnToRooms map[*websocket.Conn]*Room
}

func (rm *RoomsManager) AddRoom(roomName string) {
	rm.Rooms[roomName] = &Room{
		Name: roomName,
	}
}

func (rm *RoomsManager) EnterRoom(ws *websocket.Conn, roomName string) bool {
	if len(rm.Rooms[roomName].Players) < 2 {
		rm.ConnToRooms[ws] = rm.Rooms[roomName]
		if(rm.Rooms[roomName].Players[0].Name == "")
		{
			rm.Rooms[roomName].Players[0] = 
		}
		return true
	}

	return false
}

func (rm *RoomsManager) IsInRoom(ws *websocket.Conn) bool {
	if _, ok := rm.ConnToRooms[ws]; ok {
		return true
	}

	return false
}
