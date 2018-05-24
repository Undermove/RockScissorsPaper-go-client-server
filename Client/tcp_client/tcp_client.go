package tcp_client

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Client holds info about connection
type client struct {
	address      string // Address to open connection: localhost:9999
	localAddress string
	conn         net.Conn
	onNewMessage func(message string)
}

// Creates new tcp client instance
func New(address string) *client {
	log.Println("Connect to server with address", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client := &client{
		address: address,
		conn:    conn,
	}

	client.OnNewMessage(func(message string) {})

	return client
}

// Called when Client receives new message
func (c *client) OnNewMessage(callback func(message string)) {
	c.onNewMessage = callback
}

// Start network server
func (c *client) Listen() {
	listener, err := net.Listen("tcp", c.conn.LocalAddr().String())
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		listener.Accept()
		go c.listen()
	}
}

// Read client data from channel
func (c *client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			return
		}
		c.onNewMessage(message)
	}
}

// Send text message to client
func (c *client) Send(message string) error {
	if c.conn == nil {
		log.Println("Data not sended. Listen first!")
		return nil
	}
	_, err := c.conn.Write([]byte(message))
	return err
}
