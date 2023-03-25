package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	requestData := "request"

	requestByteData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Write(requestByteData)

	responseData := make([]byte, 100)
	n, err := conn.Read(responseData)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(responseData[:n]))
}
