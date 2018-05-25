package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type player struct {
	login string
	score int
}

func main() {
	var address string
	var port string

	var wg sync.WaitGroup

	fmt.Print("Enter server address: ")
	fmt.Fscan(os.Stdin, &address)

	fmt.Print("Enter server port: ")
	fmt.Fscan(os.Stdin, &port)

	fmt.Println("Connecting to server")

	address = "127.0.0.1"
	port = "9999"

	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected")

	player := player{score: 0}
	fmt.Print("Enter player name: ")
	fmt.Fscan(os.Stdin, &player.login)

	conn.Write([]byte("unknownUser;reg;" + player.login + ";\n"))

	wg.Add(1)
	go listen(conn, &player, &wg)

	wg.Wait()

	for {
		var command string
		fmt.Fscan(os.Stdin, &command)
		conn.Write([]byte(player.login + ";" + command + ";\n"))
	}

	io.Copy(os.Stdout, conn)
	fmt.Println("\nDone")
}

func listen(conn net.Conn, player *player, wg *sync.WaitGroup) {
	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}

		source := string(input[0:n])

		if source == "User with such name already registred." {
			fmt.Print("Enter player name: ")
			fmt.Fscan(os.Stdin, &player.login)

			conn.Write([]byte("unknownUser;reg;" + player.login + ";\n"))
		} else if source == "registred" {
			wg.Done()
		} else {
			fmt.Print("\n" + source + "\n")
		}
	}
}
