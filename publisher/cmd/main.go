package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	fmt.Fprintln(conn, "PUBLISH news HelloSubscribers!")
	fmt.Println("Published message.")
}
