package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving UDP address", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server listening on port 8080")

	buf := make([]byte, 1001)
	database := make(map[string]string)
	
  for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading data", err)
			continue
		}

		fmt.Printf("Received from %s: %s\n", clientAddr, string(buf))

		if n > 1000 {
			conn.WriteToUDP([]byte("Request size must be shorter than 1000 bytes"), clientAddr)
		}

		strContent := string(buf)

		isInsert := isInsert(strContent)

		if isInsert {
			key, value := splitKeyValue(strContent)
      fmt.Printf("Key: %v, Value: %v\n", key, value)
      database[key] = value
      continue
		}

    response := database[strContent]

		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Erro ao responder:", err)
		}
	}
}

func isInsert(strContent string) bool {
	return strings.Contains(strContent, "=")
}

func splitKeyValue(strContent string) (string, string) {
	parts := strings.SplitN(strContent, "=", 2)
	return parts[0], parts[1]
}
