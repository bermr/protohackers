package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}
		go handleConnection(conn)
		fmt.Println("Connection closed")
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	encoder := json.NewEncoder(conn)
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Println("Error receiving data", err)
			return
		}

		fmt.Println("Data received", string(data))

		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			returnError(conn)
			return
		}

		fmt.Printf("Decoded: %+v\n", req)
		if req.Method != "isPrime" {
			fmt.Println("Invalid method, closing connection")
			returnError(conn)
			return
		}

		if req.Number == nil {
			returnError(conn)
			return
		}

		fmt.Println("Request valid. Number: ", req.Number, "Method: ", req.Method)

		result := isPrime(*req.Number)
		resp := Response{
			Method: req.Method,
			Prime:  result,
		}
		err = encoder.Encode(resp)
		if err != nil {
			fmt.Println("Error sending response: ", err)
			return
		}
	}
}

func returnError(conn net.Conn) {
	response := map[string]string{
		"error": "error",
	}
	jsonBytes, _ := json.Marshal(response)
	conn.Write(append(jsonBytes, '\n'))
}

func isPrime(n float64) bool {
	if n != float64(int(n)) {
		return false
	}
	number := int(n)

	if number < 2 {
		return false
	}
	if number == 2 {
		return true
	}
	if number%2 == 0 {
		return false
	}

	for i := 3; i*i <= number; i += 2 {
		if number%i == 0 {
			return false
		}
	}
	return true
}
