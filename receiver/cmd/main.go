package main

import (
	"fmt"
	"goredis-receiver/internal/client"
	"log"
)

func main() {
	client, err := client.NewClient(&client.Options{
		Addr:     "127.0.0.1:6379",
		Username: "guest",
		Password: "guest",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	val, _ := client.Get("name")
	fmt.Println("Value:", val)

	// client.Set("password", "password")

	val, _ = client.Get("password")
	fmt.Println("Value:", val)
}
