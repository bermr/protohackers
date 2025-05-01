package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

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

	buf := make([]byte, 9)
	pricesMap := make(map[int32]int32)
	for {
		_, err := io.ReadFull(conn, buf)

		if err != nil {
			conn.Write([]byte("Error receiving data\n"))
			break
		}

		if len(buf) < 9 {
			conn.Write([]byte("Incomplete data\n"))
			break
		}

		requestType, arg1, arg2 := hexToString(buf)

		//fmt.Println("Request received", string(requestType), arg1, arg2)

		switch requestType {
		case 'I':
			pricesMap[arg1] = arg2
			/*for key, value := range pricesMap {
				fmt.Printf("Key: %d, Value: %d\n", key, value)
			}*/
			break
		case 'Q':
			mean := queryPrices(pricesMap, arg1, arg2)
			//fmt.Println("Mean: \n", mean)
			buffer := make([]byte, 4)
			binary.BigEndian.PutUint32(buffer, uint32(mean))
			conn.Write(buffer)
		default:
			conn.Write([]byte("Invalid operation\n"))
		}

	}
}

func hexToString(buf []byte) (byte, int32, int32) {
	requestType := buf[0]

	arg1 := int32(binary.BigEndian.Uint32(buf[1:5]))
	arg2 := int32(binary.BigEndian.Uint32(buf[5:9]))

	return requestType, arg1, arg2
}

func queryPrices(pricesMap map[int32]int32, minTime int32, maxTime int32) int {
	var pricesInsideInterval []int32
	var sum int = 0

	for key, value := range pricesMap {
		//fmt.Printf("Key %v, value %v, min %v, max %v \n", key, value, minTime, maxTime)
		if key >= minTime && key <= maxTime {
			pricesInsideInterval = append(pricesInsideInterval, value)
			sum += int(value)
		}
	}

  fmt.Println("Sum inside interval: ", sum)
  fmt.Println("Interval size: ", len(pricesInsideInterval))

	if sum == 0 {
		return 0
	}

	return sum / len(pricesInsideInterval)
}
