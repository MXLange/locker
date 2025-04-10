package locker

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	url string
}

// NewClient creates a new Client instance with the given URL.
// The URL should point to the server that manages the locks.
// URL must follow the format ws://serverAddr/ws?id=
// Example: ws://localhost:8080/ws?id=myapp-32
func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

// Lock establishes a WebSocket connection to the server and waits for the lock to be released.
// The connection is kept open until the lock is released.
// The function returns the WebSocket connection if successful, or nil if there was an error.
// Id the connection is lost, the process will be released.
func (c *Client) Lock(id string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s%s", c.url, id), nil)
	if err != nil {
		return nil, err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return nil, err
		}

		if strings.EqualFold(string(message), "go") {
			break
		}
	}

	return conn, nil
}

// Unlock sends a message to the server to release the lock.
func (c *Client) Unlock(conn *websocket.Conn) error {

	if conn == nil {
		return nil
	}

	return conn.WriteMessage(websocket.TextMessage, []byte("unlock"))
}
