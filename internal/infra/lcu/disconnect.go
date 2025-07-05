package lcu

func (c *Client) setDisconnected() {
	c.mu.Lock()
	changed := c.setDisconnectedLocked()
	c.mu.Unlock()
	if changed {
		c.broadcastPhase("None")
	}
}

func (c *Client) setDisconnectedLocked() bool {
	changed := c.GamePhase != "None"
	if c.Polling {
		c.log.Warn("检测到断开连接，停止轮询")
		c.stopPollingLocked()
	}
	c.Connected = false
	c.Token = ""
	c.Port = 0
	c.GamePhase = "None"
	c.failCount = 0
	return changed
}

func (c *Client) StopPolling() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stopPollingLocked()
}

func (c *Client) stopPollingLocked() {
	if !c.Polling {
		return
	}
	close(c.quit)
	c.quit = make(chan struct{})
	c.Polling = false
}
