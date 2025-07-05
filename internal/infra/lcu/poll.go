package lcu

import (
	"encoding/json"
	"net/http"
	"time"
)

const maxConsecutiveFailures = 3

func (c *Client) StartPolling() {
	c.mu.Lock()
	if c.Polling {
		c.mu.Unlock()
		return
	}
	c.Polling = true
	quit := c.quit
	c.mu.Unlock()

	c.log.Info("启动游戏状态轮询")
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.pollStatus()
			case <-quit:
				c.log.Info("停止游戏状态轮询")
				return
			}
		}
	}()
}

func (c *Client) pollStatus() {
	c.mu.RLock()
	connected := c.Connected
	c.mu.RUnlock()

	if !connected {
		return
	}

	resp, err := c.DoRequest(http.MethodGet, c.conf.Lcu.GameflowPath, nil)
	if err != nil {
		c.mu.Lock()
		c.failCount++
		if c.failCount >= maxConsecutiveFailures {
			c.log.Errorf("连续失败 %d 次，断开连接: %v", maxConsecutiveFailures, err)
			c.setDisconnectedLocked()
		} else {
			c.log.Warnf("请求失败（第 %d 次）: %v", c.failCount, err)
		}
		c.mu.Unlock()
		return
	}

	var phase string
	if err := json.Unmarshal(resp, &phase); err != nil {
		c.log.Warnf("响应解析失败: %s", string(resp))
		return
	}

	c.mu.Lock()
	c.GamePhase = phase
	c.failCount = 0
	c.mu.Unlock()

	c.log.Infof("当前游戏阶段: %s", phase)
}
