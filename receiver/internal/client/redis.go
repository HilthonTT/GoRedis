package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn    net.Conn
	scanner *bufio.Scanner
}

type Options struct {
	Addr     string
	Password string
	Username string
}

func NewClient(opt *Options) (*Client, error) {
	conn, err := net.Dial("tcp", opt.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", opt.Addr, err)
	}

	client := &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}

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
