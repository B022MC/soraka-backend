package lcu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/B022MC/soraka-backend/consts"
	"github.com/B022MC/soraka-backend/local-utils/windows/process"
	"regexp"
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
		c.log.Info("首次连接成功，开始初始化图标资源")
		itemURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ItemJson)
		champURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ChampJson)
		spellURL := fmt.Sprintf(c.conf.ProxyJsonUrl.SpellJson)
		mapURL := fmt.Sprintf(c.conf.ProxyJsonUrl.MapIconJson)
		profileURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ProfileIconJson)
		go func() {
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
		}()

	}

	c.Connected = true
	c.Port = port
	c.Token = token
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
