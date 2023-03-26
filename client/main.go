package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	RequestAndResponse(conn, "request01")
	RequestAndResponse(conn, "request02")
	RequestAndResponse(conn, "request03")
	RequestAndResponse(conn, "close")
	time.Sleep(time.Hour)

}

func RequestAndResponse(conn net.Conn, requestData string) {
	byteData, err := json.Marshal(requestData)
	if err != nil {
		return
	}

	conn.Write(byteData)
	response := make([]byte, 1024)
	num, err := conn.Read(response)
	if err != nil {
		return
	}

	fmt.Println(string(response[:num]))
}
