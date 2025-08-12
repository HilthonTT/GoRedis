package main

import (
	"fmt"
	"goredis-receiver/internal/client"
	"log"
	"time"
)

func main() {
	cli, err := client.NewClient(&client.Options{
		Addr:     "127.0.0.1:6379",
		Username: "guest",
		Password: "guest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 1. Test Set
	fmt.Println("Setting name=John...")
	if err := cli.Set("name", "John"); err != nil {
		log.Fatal("SET failed:", err)
	}

	// 2. Test Get
	val, err := cli.Get("name")
	if err != nil {
		log.Fatal("GET failed:", err)
	}
	fmt.Println("GET name:", val)

	// 3. Test Subscribe
	fmt.Println("Subscribing to topic 'news'...")

	done := make(chan struct{})

	err = cli.Subscribe("news", func(msg string) {
		fmt.Println("Got message from subscription:", msg)
	}, func(err error) {
		fmt.Println("Subscription error:", err)
		close(done)
	})

	if err != nil {
		log.Fatal("SUBSCRIBE failed:", err)
	}

	// 4. Simulate publishing from the same process after a short delay
	go func() {
		time.Sleep(2 * time.Second)
		cli.Publish("news", "BreakingNews!")
	}()

	// Wait here until the subscription signals it has ended
	<-done
	fmt.Println("Subscription ended, exiting")
}
