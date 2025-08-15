package main

import (
	"fmt"
	"goredis-shared/redis"
	"log"
	"time"
)

func main() {
	cli, err := redis.NewClient(&redis.Options{
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

	// err = cli.Subscribe("news", func(msg string) {
	// 	fmt.Println("Got message from subscription:", msg)
	// }, func(err error) {
	// 	fmt.Println("Subscription error:", err)
	// close(done)
	// })

	if err != nil {
		log.Fatal("SUBSCRIBE failed:", err)
	}

	// 4. Simulate publishing from the same process after a short delay
	go func() {
		time.Sleep(2 * time.Second)
		cli.Publish("news", "BreakingNews!")
	}()

	// 4. Test SADD
	fmt.Println("Adding fruits to 'myset'...")
	if err := cli.SAdd("myset", "apple", "banana", "cherry"); err != nil {
		log.Fatal("SADD failed:", err)
	}

	// 5. Test SMEMBERS
	fmt.Println("Getting members of 'myset'...")
	members, err := cli.SMembers("myset")
	if err != nil {
		log.Fatal("SMEMBERS failed:", err)
	}
	fmt.Println("Members:", members)

	// Wait here until the subscription signals it has ended
	close(done)
	<-done
	fmt.Println("Subscription ended, exiting")
}
