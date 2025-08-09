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

	fmt.Fprintln(conn, "GET name")

	// fmt.Fprintln(conn, "SET name John")
	// fmt.Fprintln(conn, "GET name")

	// fmt.Fprintln(conn, "DEL name")
	// fmt.Fprintln(conn, "GET name")

	// Wrong commands
	// fmt.Fprintln(conn, "GET")
	// fmt.Fprintln(conn, "SET name")
	// fmt.Fprintln(conn, "DEL name test")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Server:", scanner.Text())
	}
}
