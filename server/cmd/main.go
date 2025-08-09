package main

import (
	"bufio"
	"fmt"
	"goredis-server/internal/cache"
	"goredis-server/internal/data"
	"goredis-server/internal/messaging"
	"net"
	"os"
	"strings"
)

var db = cache.NewShardMap(16)

func main() {
	loadSnapshot()

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
			conn.Write([]byte("OK\n"))
			data.LogCommand("SET", key, value)
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
			data.LogCommand("DEL", key, "")
		case "SUBSCRIBE":
			messaging.HandleSubscribe(conn, args[1])
			conn.Write([]byte("OK\n"))
		case "PUBLISH":
			topic := args[1]
			message := args[2]
			messaging.HandlePublish(topic, message)
			conn.Write([]byte("OK\n"))
		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}

	data.LogFile.Close()
}

func loadSnapshot() {
	file, err := os.Open("snap.log")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd := strings.Fields(line)

		cmdUpper := strings.ToUpper(cmd[0])
		key := cmd[1]
		switch cmdUpper {
		case "SET":
			if len(cmd) < 3 {
				continue
			}

			value := cmd[2]
			db.Set(key, value)
		case "DEL":
			db.Delete(key)
		}
	}
}
