package redis

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn    net.Conn
	scanner *bufio.Scanner
	mu      sync.RWMutex
}

type Options struct {
	Addr     string
	Password string
	Username string
}

func NewClient(opt *Options) (*Client, error) {
	conn, err := net.DialTimeout("tcp", opt.Addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", opt.Addr, err)
	}

	client := &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}

	buf := make([]byte, 0, 64*1024)
	client.scanner.Buffer(buf, 1024*1024)

	// AUTH command (fmt.Fprintln adds newline automatically)
	authCmd := fmt.Sprintf("AUTH %s %s", opt.Username, opt.Password)
	if _, err := fmt.Fprintln(client.conn, authCmd); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to send AUTH command: %w", err)
	}

	resp, err := client.readResponse()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to read AUTH response: %w", err)
	}
	if resp != "OK" {
		client.Close()
		return nil, fmt.Errorf("authentication failed: %s", resp)
	}

	return client, nil
}

func (c *Client) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cmd := fmt.Sprintf("GET %s", key)
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return "", fmt.Errorf("failed to send GET command: %w", err)
	}

	resp, err := c.readResponse()
	if err != nil {
		return "", fmt.Errorf("failed to read GET response: %w", err)
	}

	if resp == "(nil)" {
		return "", nil // key not found
	}

	return resp, nil
}

func (c *Client) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := fmt.Sprintf("SET %s %s", key, value)
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return fmt.Errorf("failed to send SET command: %w", err)
	}

	resp, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read SET response: %w", err)
	}

	if resp != "OK" {
		return fmt.Errorf("SET failed: %s", resp)
	}

	return nil
}

func (c *Client) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := fmt.Sprintf("DEL %s", key)
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return fmt.Errorf("failed to send DEL command: %w", err)
	}

	return nil
}

func (c *Client) SRem(key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	cmd := "SREM " + key
	for _, m := range members {
		cmd += " " + m
	}

	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return fmt.Errorf("failed to send SREM command: %w", err)
	}

	_, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read SREM response: %w", err)
	}

	return nil
}

func (c *Client) Publish(channel, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := fmt.Sprintf("PUBLISH %s %s", channel, value)
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return fmt.Errorf("failed to send PUBLISH command: %w", err)
	}

	resp, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read PUBLISH response: %w", err)
	}

	if resp != "OK" {
		return fmt.Errorf("PUBLISH failed: %s", resp)
	}

	return nil
}

func (c *Client) Subscribe(channel string, onMessage func(msg string), onError func(error)) error {
	if strings.TrimSpace(channel) == "" {
		return fmt.Errorf("channel name cannot be empty")
	}

	c.mu.Lock()
	_, err := fmt.Fprintf(c.conn, "SUBSCRIBE %s\n", channel)
	c.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to send SUBSCRIBE command: %w", err)
	}

	go func() {
		for c.scanner.Scan() {
			msg := strings.TrimSpace(c.scanner.Text())
			if msg != "" && msg != "OK" {
				onMessage(msg)
			}
		}

		if err := c.scanner.Err(); err != nil {
			onError(fmt.Errorf("subscribe listen error: %w", err))
		} else {
			onError(nil)
		}
	}()
	return nil
}

func (c *Client) SAdd(key string, members ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := fmt.Sprintf("SADD %s %s", key, strings.Join(members, " "))
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return fmt.Errorf("failed to send SADD command: %w", err)
	}

	resp, err := c.readResponse()
	if err != nil {
		return fmt.Errorf("failed to read SADD response: %w", err)
	}

	if resp == "(error)" {
		return fmt.Errorf("SADD failed: %s", resp)
	}
	return nil
}

func (c *Client) SMembers(key string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmd := fmt.Sprintf("SMEMBERS %s", key)
	if _, err := fmt.Fprintln(c.conn, cmd); err != nil {
		return nil, fmt.Errorf("failed to send SMEMBERS command: %w", err)
	}

	line, err := c.readResponse()
	if err != nil {
		return nil, err
	}

	var n int
	if _, err := fmt.Sscanf(line, "*%d", &n); err != nil {
		return nil, fmt.Errorf("invalid array header: %s", line)
	}

	members := make([]string, 0, n)
	for i := 0; i < n; i++ {
		// Read bulk string length (skip it)
		_, err := c.readResponse()
		if err != nil {
			return nil, err
		}

		// Read actual value
		val, err := c.readResponse()
		if err != nil {
			return nil, err
		}
		members = append(members, val)
	}

	return members, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) readResponse() (string, error) {
	if c.scanner.Scan() {
		return strings.TrimSpace(c.scanner.Text()), nil
	}
	if err := c.scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("connection closed by server")
}
