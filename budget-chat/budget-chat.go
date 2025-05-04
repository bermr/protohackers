package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080")
	clients := make(map[string]net.Conn)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}

		go handleConnection(conn, clients)
	}
}

func handleConnection(conn net.Conn, clients map[string]net.Conn) {
	clientName := ""

	defer func() {
		if clientName != "" {
			broadcastButNotTo(fmt.Sprintf("* %s já ralou peito\n", clientName), clientName, clients)
			delete(clients, clientName)
		}
		conn.Close()
		fmt.Println("Connection closed", conn.RemoteAddr())
	}()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	conn.Write([]byte("Aoba. Cumé q c chama?\n"))
	fmt.Println("Waiting for client name")

	line, err := reader.ReadString('\n')

	if err != nil {
		conn.Write([]byte("Error receiving data\n"))
		return
	}

	line = strings.TrimSpace(line)

	if len(line) < 1 {
		conn.Write([]byte("Error: client name is empty\n"))
		fmt.Println("Client name is empty")
		return
	}

	for _, r := range line {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			fmt.Println("Invalid character in client name", r)
			conn.Write([]byte("Error: client name contains non-ASCII characters\n"))
			return
		}
	}

	fmt.Println("Client name", line)
	clientName = line
	clients[clientName] = conn

	var keys []string
	for key := range clients {
		if key != clientName {
			keys = append(keys, key)
		}
	}
	newClientMessage := fmt.Sprint("* A galera que tá ai: ", strings.Join(keys, ", "), "\n")
	fmt.Print(newClientMessage)
	sendMessageTo(newClientMessage, clientName, clients)

	welcomeMessage := fmt.Sprintf("* %s entrou no chat\n", clientName)
	fmt.Print(welcomeMessage)
	broadcastButNotTo(welcomeMessage, clientName, clients)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			conn.Write([]byte("Error receiving data\n"))
			break
		}

		msg := fmt.Sprintf("[%s] %s", clientName, line)
		fmt.Print(msg)
		broadcastButNotTo(msg, clientName, clients)
	}
}

func sendMessageTo(message string, clientName string, clients map[string]net.Conn) {
	for key, conn := range clients {
		if clientName == key {
			conn.Write([]byte(message))
		}
	}
}

func broadcastButNotTo(message string, clientName string, clients map[string]net.Conn) {
	if len(clients) == 0 {
		fmt.Println("No clients to broadcast to")
		return
	}

	for key, conn := range clients {
		if clientName != key {
			conn.Write([]byte(message))
		}
	}
}
