package messaging

import (
	"log"
	"net"
	"sync"
)

type subscriber struct {
	conn net.Conn
	ch   chan string
}

type PubSub struct {
	mu     sync.RWMutex
	topics map[string][]*subscriber
}

var pubsub = PubSub{
	topics: make(map[string][]*subscriber),
}

func HandleSubscribe(conn net.Conn, topic string) {
	ch := make(chan string, 16)

	sub := &subscriber{
		conn: conn,
		ch:   ch,
	}

	pubsub.mu.Lock()
	pubsub.topics[topic] = append(pubsub.topics[topic], sub)
	pubsub.mu.Unlock()

	log.Printf("Client %v subscribed to %q\n", conn.RemoteAddr(), topic)

	go func() {
		defer func() {
			pubsub.removeSubscriber(topic, sub)
			conn.Close()
			log.Printf("Client %v unsubscribed from %q\n", conn.RemoteAddr(), topic)
		}()

		for msg := range ch {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				log.Printf("Error writing to subscriber: %v: %v", conn.RemoteAddr(), err)
				return
			}
		}
	}()
}

func HandlePublish(topic, message string) {
	pubsub.mu.RLock()
	defer pubsub.mu.RUnlock()

	subs := pubsub.topics[topic]
	if len(subs) == 0 {
		log.Printf("No subscribers for topic %q\n", topic)
		return
	}

	log.Printf("Publishing to %d subscriber(s) on %q\n", len(subs), topic)

	for _, sub := range subs {
		select {
		case sub.ch <- message:
			// sent successfully
		default:
			log.Printf("Subscriber %v is too slow, skipping message\n", sub.conn.RemoteAddr())
		}
	}
}

func (ps *PubSub) removeSubscriber(topic string, sub *subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs := ps.topics[topic]
	for i, s := range subs {
		if s == sub {
			ps.topics[topic] = append(subs[:i], subs[i+1:]...)
			close(s.ch)
			break
		}
	}
}
