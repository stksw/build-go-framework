package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listener.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	handle(conn)
}

func handle(conn net.Conn) {

	buf := make([]byte, 1000)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("request:")
	fmt.Println(string(buf[:n]))

	responseData := "response"

	responseByteData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(responseByteData)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Close()
}
