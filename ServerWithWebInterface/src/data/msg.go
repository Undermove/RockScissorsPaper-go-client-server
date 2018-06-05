package msg

import (
	"encoding/json"
)

// Define our message object
type Message struct {
	Type string
	Raw  json.RawMessage
}

// Define our message object
type AuthRequest struct {
	Username string `json:"username"`
}

// Define our message object
type AuthResponse struct {
	IsRegistred  bool   `json:"isRegistred"`
	RejectReason string `json:"rejectReason"`
}

// Define our message object
type CreateRoomRequest struct {
	RoomName string `json:"rejectReason"`
}

// Define our message object
type CreateRoomResponse struct {
	IsCreated    bool   `json:"isRegistred"`
	RejectReason string `json:"rejectReason"`
}
