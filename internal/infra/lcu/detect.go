package lcu

import (
	"fmt"
	"strconv"
	"time"

	"regexp"

	"github.com/B022MC/soraka-backend/consts"
	"github.com/B022MC/soraka-backend/local-utils/windows/process"
)

var lolCommandlineReg = regexp.MustCompile(`--remoting-auth-token=(.+?)" ".*?--app-port=(\d+)"`)

func (c *Client) backgroundLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.tryDetectClient()
		case <-c.quit:
			c.log.Info("检测循环退出")
			return
		}
	}
}

func (c *Client) tryDetectClient() {
	port, token, err := getLolClientApiInfo()

	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		if c.Connected {
			c.log.Warn("LCU 客户端已关闭")
		}
		c.setDisconnectedLocked()
		return
	}

	if !c.Connected && !c.Polling {
		c.log.Infof("检测到 LCU，端口 %d", port)
		go c.StartPolling()
		//go c.initResourcesIfNeeded()

		// 异步等待LCU API准备好后再初始化资源
		go func(p int, t string) {
			// 等待LCU API可用（最多重试10次，每次间隔1秒）
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second)
				if c.waitForLCUReady(p, t) {
					break
				}
				c.log.Infof("等待LCU API准备就绪... (%d/10)", i+1)
			}

			c.log.Info("首次连接成功，开始初始化图标资源")
			itemURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ItemJson)
			champURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ChampJson)
			spellURL := fmt.Sprintf(c.conf.ProxyJsonUrl.SpellJson)
			mapURL := fmt.Sprintf(c.conf.ProxyJsonUrl.MapIconJson)
			profileURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ProfileIconJson)
			err := c.IconMapDownloader(
				itemURL,
				champURL,
				spellURL,
				profileURL,
				mapURL,
			)
			if err != nil {
				c.log.Warnf("图标资源初始化失败: %v", err)
			} else {
				c.log.Info("图标资源初始化完成")
			}
		}(port, token)

	}

	c.Connected = true
	c.Port = port
	c.Token = token

	// 启动 WebSocket 连接
	if c.wsManager != nil && !c.Polling {
		go func() {
			if err := c.wsManager.Start(); err != nil {
				c.log.Errorf("Failed to start WebSocket: %v", err)
			}
		}()
	}
}

func getLolClientApiInfo() (port int, token string, err error) {
	cmdline, err := process.GetProcessCommand(consts.LolUxProcessName)
	if err != nil {
		return 0, "", fmt.Errorf("获取进程失败: %w", err)
	}

	matches := lolCommandlineReg.FindSubmatch([]byte(cmdline))
	if len(matches) < 3 {
		return 0, "", fmt.Errorf("命令行参数不完整")
	}
	token = string(matches[1])
	port, err = strconv.Atoi(string(matches[2]))
	return
}

// waitForLCUReady 检测LCU API是否准备就绪
func (c *Client) waitForLCUReady(port int, token string) bool {
	// 临时设置端口和token来尝试请求
	c.mu.Lock()
	c.Port = port
	c.Token = token
	c.mu.Unlock()

	// 尝试一个简单的API请求来检测是否准备好
	_, err := c.DoRequest("GET", "/lol-summoner/v1/current-summoner", nil)
	return err == nil
}
