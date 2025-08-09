package main

import (
	"bufio"
	"fmt"
	"goredis-server/internal/cache"
	"net"
	"strings"
)

var db = cache.NewShardMap(16)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started on port 6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "SET":
			if len(args) != 3 {
				fmt.Fprintln(conn, "ERR wrong arguments")
				continue
			}

			key, value := args[1], args[2]
			db.Set(key, value)
		case "GET":
			if len(args) != 2 {
				fmt.Fprintln(conn, "ERR wrong arguments")
				continue
			}

			key := args[1]
			val, ok := db.Get(key)
			if !ok {
				fmt.Fprintln(conn, "(nil)")
			} else {
				fmt.Fprintln(conn, val)
			}
		case "DEL":
			if len(args) != 2 {
				fmt.Fprintln(conn, "ERR wrong arguments")
				continue
			}

			key := args[1]
			db.Delete(key)
		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}
}
