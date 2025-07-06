package lcu

import "sync"

func (c *Client) Mutex() *sync.RWMutex {
	return &c.mu
}

func (c *Client) IsConnected() bool {
	return c.Connected
}

func (c *Client) GetGamePhase() string {
	return c.GamePhase
}

func (c *Client) GetToken() string {
	return c.Token
}

func (c *Client) GetPort() int {
	return c.Port
}

func (c *Client) GetClientPath() string {
	return c.ClientPath
}
