package lcu

import "github.com/B022MC/soraka-backend/consts"

func (c *Client) setDisconnected() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setDisconnectedLocked()
}

func (c *Client) setDisconnectedLocked() {
	if c.Polling {
		c.log.Warn("检测到断开连接，停止轮询")
		c.stopPollingLocked()
	}
	c.Connected = false
	c.Token = ""
	c.Port = 0
	c.GamePhase = "None"
	c.failCount = 0
	// 清空图标缓存
	for k := range consts.ItemIconMap {
		delete(consts.ItemIconMap, k)
	}
	for k := range consts.SpellIconMap {
		delete(consts.SpellIconMap, k)
	}
	for k := range consts.ChampIconMap {
		delete(consts.ChampIconMap, k)
	}
	for k := range consts.MapIcon {
		delete(consts.MapIcon, k)
	}
	for k := range consts.ProfileIconMap {
		delete(consts.ProfileIconMap, k)
	}
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
