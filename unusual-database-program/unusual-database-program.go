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

	database := make(map[string]string)

	for {
		buf := make([]byte, 1025)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading data", err)
			continue
		}

		fmt.Printf("Received from %s: %s\n", addr, string(buf[:n]))

		if n > 1024 {
			fmt.Println("Request too large")
			conn.WriteTo([]byte("Request size must be shorter than 1000 bytes"), addr)
			continue
		}

		strContent := string(buf[:n])
		strContent = strings.TrimRight(strContent, " \t\r\n")

		isInsert := isInsert(strContent)

		if isInsert {
			key, value := splitKeyValue(strContent)
			fmt.Printf("Key: %v, Value: %v\n", key, value)
			if key != "version" {
				database[key] = value
			}
			conn.WriteTo([]byte(""), addr)
			continue
		}

		var response string
		if len(database) > 0 {
			_, ok := database[strContent]

			if ok {
				response = fmt.Sprintf("%v=%v", strContent, database[strContent])
			}
		}

		if strContent == "version" {
			response = "version=Unusual Database do B"
		}

		_, err = conn.WriteTo([]byte(response), addr)
		if err != nil {
			fmt.Println("Response error", err)
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
