package lcu

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) initResourcesIfNeeded() {
	c.initOnce.Do(func() {
		path, err := findLolPath()
		if err != nil {
			c.log.Warnf("自动查找 LOL 路径失败: %v", err)
		} else {
			c.ClientPath = path
			c.log.Infof("已自动记录 LOL 安装路径: %s", path)
		}
		c.log.Info("首次连接成功，开始初始化图标资源")
		itemURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ItemJson)
		champURL := fmt.Sprintf(c.conf.ProxyJsonUrl.ChampJson)
		spellURL := fmt.Sprintf(c.conf.ProxyJsonUrl.SpellJson)
		err = c.IconMapDownloader(
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
func findLolPath() (string, error) {
	basePath := getLoLPathByRegistry()
	if basePath != "" {
		if exePath := appendClientExe(basePath); exePath != "" {
			return exePath, nil
		}
	}

	basePath, err := findLolPathFromUninstall()
	if err != nil {
		return "", err
	}
	if exePath := appendClientExe(basePath); exePath != "" {
		return exePath, nil
	}

	return "", fmt.Errorf("未找到 LOL 客户端可执行文件")
}

// getLoLPathByRegistry 通过国服路径尝试查找 LOL 安装目录
func getLoLPathByRegistry() string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Tencent\LOL`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	val, _, err := key.GetStringValue("InstallPath")
	if err != nil || val == "" {
		return ""
	}

	path := filepath.Join(val, "TCLS")
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, `\`, `/`)

	if len(path) > 1 {
		path = strings.ToUpper(path[:1]) + path[1:]
	}
	return path
}

// findLolPathFromUninstall 从常见的卸载表注册表中查找 LOL 安装位置
func findLolPathFromUninstall() (string, error) {
	roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}
	subkeys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}
	for _, root := range roots {
		for _, sub := range subkeys {
			k, err := registry.OpenKey(root, sub, registry.READ)
			if err != nil {
				continue
			}
			names, err := k.ReadSubKeyNames(-1)
			if err != nil {
				k.Close()
				continue
			}
			for _, name := range names {
				sk, err := registry.OpenKey(k, name, registry.READ)
				if err != nil {
					continue
				}
				display, _, _ := sk.GetStringValue("DisplayName")
				if strings.Contains(display, "英雄联盟") || strings.Contains(display, "League of Legends") {
					if loc, _, err := sk.GetStringValue("InstallLocation"); err == nil && loc != "" {
						sk.Close()
						k.Close()
						return filepath.Join(loc, "League of Legends.exe"), nil
					}
					if icon, _, err := sk.GetStringValue("DisplayIcon"); err == nil && icon != "" {
						sk.Close()
						k.Close()
						return icon, nil
					}
				}
				sk.Close()
			}
			k.Close()
		}
	}
	return "", fmt.Errorf("LOL path not found in registry")
}

// appendClientExe 尝试在路径上拼接常见可执行文件名，并检查存在性
func appendClientExe(basePath string) string {
	candidates := []string{
		filepath.Join(basePath, "LeagueClient.exe"),
		filepath.Join(basePath, "Client.exe"),
	}

	for _, path := range candidates {
		if fileExists(path) {
			return path
		}
	}
	return ""
}

// fileExists 检查文件是否存在且不是目录
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
