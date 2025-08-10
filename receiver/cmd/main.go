package main

import (
	"bufio"
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

	// Send AUTH before anything else
	fmt.Fprintln(conn, "AUTH guest guest")

	// Then send your actual command
	fmt.Fprintln(conn, "SUBSCRIBE news")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Message from server:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}
