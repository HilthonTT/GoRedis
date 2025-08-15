package main

import (
	"bufio"
	"fmt"
	"goredis-server/internal/cache"
	"goredis-server/internal/config"
	"goredis-server/internal/data"
	"goredis-server/internal/handler"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	cfg := config.NewConfig()

	handler := handler.NewHandler()

	loadSnapshot(handler.DB)
	defer data.CloseLog()

	ln, err := net.Listen("tcp", net.JoinHostPort(cfg.BindAddr, strconv.Itoa(cfg.Port)))
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Printf("Server started on port %d\n", cfg.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		handler.SetConn(conn)
		go handleConnection(conn, handler, cfg)
	}
}

func handleConnection(conn net.Conn, handler *handler.Handler, cfg *config.Config) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	authenticated := false

	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		cmd := strings.ToUpper(args[0])

		if !authenticated {
			if strings.ToUpper(args[0]) == "AUTH" {
				authenticated = handler.Auth(args, cfg)
				if !authenticated {
					fmt.Fprintln(conn, "ERR invalid username or password")
					return // close connection
				}
			} else {
				fmt.Fprintln(conn, "NOAUTH Authentication required.")
			}
			continue
		}

		fmt.Printf("CMD: %v\n", cmd)

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
		case "SADD":
			handler.SAdd(args)
		case "SMEMBERS":
			handler.SMembers(args)
		case "SREM":
			handler.SRem(args)
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

		// Split into [timestamp, command...]
		parts := strings.SplitN(line, " | ", 2)
		if len(parts) != 2 {
			continue // malformed line
		}

		// Parse command from second part
		cmd := strings.Fields(parts[1])
		if len(cmd) < 2 {
			continue // invalid command
		}

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

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading snapshot: %v\n", err)
	}
}
