package lcu

import "fmt"

func (c *Client) initResourcesIfNeeded() {
	c.initOnce.Do(func() {
		c.log.Info("首次连接成功，开始初始化图标资源")
		itemURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ItemJson)
		champURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ChampJson)
		spellURL := fmt.Sprintf(c.conf.ProxyJsonUrl.SpellJson)
		err := c.IconMapDownloader(
			itemURL,
			champURL,
			spellURL,
		)
		if err != nil {
			c.log.Warnf("图标资源初始化失败: %v", err)
		} else {
			c.log.Info("图标资源初始化完成")
		}
	})
}
