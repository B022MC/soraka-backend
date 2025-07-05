package lcu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IconEntry struct {
	ID       int    `json:"id"`
	IconPath string `json:"iconPath"`
}

var (
	ItemIconMap  = make(map[int]string)
	ChampIconMap = make(map[int]string)
	SpellIconMap = make(map[int]string)
)

// 自动下载并初始化所有 icon map
func (c *Client) IconMapDownloader(itemURL, champURL, spellURL string) error {
	if err := c.downloadAndLoad(itemURL, ItemIconMap, "item"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(champURL, ChampIconMap, "champion"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(spellURL, SpellIconMap, "spell"); err != nil {
		return err
	}
	fmt.Println("所有图标 map 初始化完成")
	return nil
}

func (c *Client) downloadAndLoad(url string, target map[int]string, logName string) error {
	fmt.Printf("[%s] 正在下载: %s\n", logName, url)
	if url == "" {
		return nil
	}
	resp, err := c.DoRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("[%s] 下载失败: %w", logName, err)
	}
	var entries []IconEntry
	if err := json.Unmarshal(resp, &entries); err != nil {
		return fmt.Errorf("[%s] 解析 JSON 失败: %w", logName, err)
	}

	for _, e := range entries {
		target[e.ID] = e.IconPath
	}

	fmt.Printf("[%s] 加载 %d 个图标\n", logName, len(target))
	return nil
}
