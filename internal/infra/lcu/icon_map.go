package lcu

import (
	"encoding/json"
	"fmt"
	"github.com/B022MC/soraka-backend/consts"
	"net/http"
)

type IconEntry struct {
	ID       int    `json:"id"`
	IconPath string `json:"iconPath"`
}
type SpellIconEntry struct {
	ID       int    `json:"id"`
	IconPath string `json:"iconPath"`
}

type ItemIconEntry struct {
	ID       int    `json:"id"`
	IconPath string `json:"iconPath"`
}

type ChampIconEntry struct {
	ID                 int    `json:"id"`
	SquarePortraitPath string `json:"squarePortraitPath"`
}
type MapIconEntry struct {
	ID          int    `json:"id"`
	MapStringId string `json:"mapStringId"`
}

type ProfileIconEntry struct {
	ID       int    `json:"id"`
	IconPath string `json:"iconPath"`
}

// 自动下载并初始化所有 icon map
func (c *Client) IconMapDownloader(itemURL, champURL, spellURL, profileURL, mapURL string) error {
	if err := c.downloadAndLoad(itemURL, consts.ItemIconMap, "item"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(champURL, consts.ChampIconMap, "champion"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(spellURL, consts.SpellIconMap, "spell"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(mapURL, consts.MapIcon, "map"); err != nil {
		return err
	}
	if err := c.downloadAndLoad(profileURL, consts.ProfileIconMap, "profile"); err != nil {
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

	switch logName {
	case "item", "spell":
		var entries []SpellIconEntry
		if err := json.Unmarshal(resp, &entries); err != nil {
			return fmt.Errorf("[%s] 解析 JSON 失败: %w", logName, err)
		}
		for _, e := range entries {
			target[e.ID] = e.IconPath
		}
	case "champion":
		var entries []ChampIconEntry
		if err := json.Unmarshal(resp, &entries); err != nil {
			return fmt.Errorf("[%s] 解析 JSON 失败: %w", logName, err)
		}
		for _, e := range entries {
			target[e.ID] = e.SquarePortraitPath
		}
	case "map":
		var entries []MapIconEntry
		if err := json.Unmarshal(resp, &entries); err != nil {
			return fmt.Errorf("[%s] 解析 JSON 失败: %w", logName, err)
		}
		for _, e := range entries {
			target[e.ID] = e.MapStringId
		}
	case "profile":
		var entries []ProfileIconEntry
		if err := json.Unmarshal(resp, &entries); err != nil {
			return fmt.Errorf("[%s] 解析 JSON 失败: %w", logName, err)
		}
		for _, e := range entries {
			target[e.ID] = e.IconPath
		}
	default:
		return fmt.Errorf("[%s] 未知类型，无法解析", logName)
	}

	fmt.Printf("[%s] 加载 %d 个图标\n", logName, len(target))
	return nil
}
