package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type Request struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
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
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	var buf bytes.Buffer
	tee := io.TeeReader(conn, &buf)
	decoder := json.NewDecoder(tee)
	encoder := json.NewEncoder(conn)

	for {

		var req Request
		err := decoder.Decode(&req)
		fmt.Println("Data received", buf.String())
		buf.Reset()
		
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed.")
			} else {
				fmt.Println("Error decoding JSON", err)
			}
			conn.Close()
			break
		}

		if req.Method != "isPrime" {
			fmt.Println("Invalid method, closing connection")
			break
		}

		if req.Number == 0 {
			break
		}

		fmt.Println("Request valid. Number: ", req.Number, "Method: ", req.Method)

		result := isPrime(req.Number)
		resp := Response{
			Method: req.Method,
			Prime:  result,
		}
		err = encoder.Encode(resp)
		if err != nil {
			fmt.Println("Error sending response: ", err)
			break
		}
	}
	fmt.Println("Done")
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
