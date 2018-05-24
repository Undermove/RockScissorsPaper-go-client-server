package main

import (
	"fmt"
	"net"
	"strings"
)

var dict = map[string]int{
	"register": 1,
	"turn":     2,
}

type player struct {
	login   string
	score   int
	address string
	port    string
}

func main() {
	var players []player
	var connections []net.Conn
	var isGameStarted = false

	fmt.Print("Starting Server...")

	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	fmt.Println("Server is listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		connections = append(connections, conn)
		go handleRegistration(conn, &players, &isGameStarted, &connections)
	}
}

// обработка подключения
func handleRegistration(conn net.Conn, players *[]player, isGameStarted *bool, connections *[]net.Conn) {
	defer conn.Close()
	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		source := string(input[0:n])

		command := strings.Split(source, ":")

		target, ok := dict[command[0]]
		if ok == false { // если данные не найдены в словаре
			conn.Write([]byte("Wrong command"))
			break
		}

		if target == 1 {
			newPlayer := player{
				login:   command[1],
				score:   0,
				address: strings.Split(conn.RemoteAddr().String(), ":")[0],
				port:    strings.Split(conn.RemoteAddr().String(), ":")[1],
			}

			*players = append(*players, newPlayer)
			// выводим на консоль сервера диагностическую информацию
			fmt.Println(source, "-", target)
			// отправляем данные клиенту
			conn.Write([]byte(newPlayer.login + " registred"))
		} else if target == 2 && *isGameStarted == true {
			fmt.Println("Turn!!!")
			// отправляем данные клиенту
			conn.Write([]byte("Turn!"))
			continue
		}

		playersCount := len(*players)

		if playersCount == 2 {
			for _, value := range *connections {
				value.Write([]byte("Roomisfull"))
			}

			*isGameStarted = true
		}

		for _, value := range *players {
			fmt.Println(value.login)
		}
	}
}
