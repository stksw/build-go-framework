package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	for {
		time.Sleep(time.Second * 3)
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		requestInfo := string(buf[:n])
		fmt.Println("request:")
		fmt.Println(requestInfo)

		if requestInfo == `"close"` {
			fmt.Println("close from connecting")
			conn.Close()
			return
		}

		response := "response"
		responseByteData, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = conn.Write(responseByteData)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
