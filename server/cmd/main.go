package main

import (
	"bufio"
	"fmt"
	"goredis-server/internal/cache"
	"goredis-server/internal/data"
	"goredis-server/internal/handler"
	"net"
	"os"
	"strings"
)

func main() {
	handler := handler.NewHandler()

	loadSnapshot(handler.DB)
	defer data.CloseLog()

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
		handler.SetConn(conn)
		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler *handler.Handler) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "SET":
			handler.Set(args)
		case "GET":
			handler.Get(args)
		case "SUBSCRIBE":
			handler.Subscribe(args)
		case "PUBLISH":
			handler.Publish(args)
		case "EXPIRE":
			handler.Expire(args)
		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}
}

func loadSnapshot(db cache.ShardMap) {
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
