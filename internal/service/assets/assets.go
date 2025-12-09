package assets

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	DDRAGON_VERSION = "14.23.1"
	DDRAGON_CDN     = "https://ddragon.leagueoflegends.com/cdn/" + DDRAGON_VERSION
	CDRAGON_CDN     = "https://raw.communitydragon.org/latest"
)

type AssetsService struct {
	cacheDir string
	mu       sync.RWMutex
}

func NewAssetsService() *AssetsService {
	// 获取可执行文件所在目录
	exePath, _ := os.Executable()
	cacheDir := filepath.Join(filepath.Dir(exePath), "assets_cache")

	// 创建缓存目录
	os.MkdirAll(filepath.Join(cacheDir, "champion-icons"), 0755)
	os.MkdirAll(filepath.Join(cacheDir, "profile-icons"), 0755)
	os.MkdirAll(filepath.Join(cacheDir, "items"), 0755)
	os.MkdirAll(filepath.Join(cacheDir, "spells"), 0755)
	os.MkdirAll(filepath.Join(cacheDir, "perks"), 0755)

	return &AssetsService{
		cacheDir: cacheDir,
	}
}

// GetChampionIcon 获取英雄图标
func (s *AssetsService) GetChampionIcon(championId int) ([]byte, string, error) {
	filename := fmt.Sprintf("%d.png", championId)
	cachePath := filepath.Join(s.cacheDir, "champion-icons", filename)

	// 检查缓存
	if data, err := s.readCache(cachePath); err == nil {
		return data, "image/png", nil
	}

	// 从CDN下载
	url := fmt.Sprintf("%s/plugins/rcp-be-lol-game-data/global/default/v1/champion-icons/%d.png", CDRAGON_CDN, championId)
	data, err := s.downloadAndCache(url, cachePath)
	if err != nil {
		return nil, "", err
	}

	return data, "image/png", nil
}

// GetProfileIcon 获取召唤师头像
func (s *AssetsService) GetProfileIcon(iconId int) ([]byte, string, error) {
	filename := fmt.Sprintf("%d.jpg", iconId)
	cachePath := filepath.Join(s.cacheDir, "profile-icons", filename)

	// 检查缓存
	if data, err := s.readCache(cachePath); err == nil {
		return data, "image/jpeg", nil
	}

	// 从CDN下载
	url := fmt.Sprintf("%s/plugins/rcp-be-lol-game-data/global/default/v1/profile-icons/%d.jpg", CDRAGON_CDN, iconId)
	data, err := s.downloadAndCache(url, cachePath)
	if err != nil {
		return nil, "", err
	}

	return data, "image/jpeg", nil
}

// GetItemIcon 获取装备图标
func (s *AssetsService) GetItemIcon(itemId int) ([]byte, string, error) {
	filename := fmt.Sprintf("%d.png", itemId)
	cachePath := filepath.Join(s.cacheDir, "items", filename)

	// 检查缓存
	if data, err := s.readCache(cachePath); err == nil {
		return data, "image/png", nil
	}

	// 从CDN下载
	url := fmt.Sprintf("%s/img/item/%d.png", DDRAGON_CDN, itemId)
	data, err := s.downloadAndCache(url, cachePath)
	if err != nil {
		return nil, "", err
	}

	return data, "image/png", nil
}

// GetSpellIcon 获取召唤师技能图标
func (s *AssetsService) GetSpellIcon(spellName string) ([]byte, string, error) {
	filename := fmt.Sprintf("%s.png", spellName)
	cachePath := filepath.Join(s.cacheDir, "spells", filename)

	// 检查缓存
	if data, err := s.readCache(cachePath); err == nil {
		return data, "image/png", nil
	}

	// 从CDN下载
	url := fmt.Sprintf("%s/img/spell/%s.png", DDRAGON_CDN, spellName)
	data, err := s.downloadAndCache(url, cachePath)
	if err != nil {
		return nil, "", err
	}

	return data, "image/png", nil
}

// GetPerkIcon 获取符文图标
func (s *AssetsService) GetPerkIcon(perkStyle string) ([]byte, string, error) {
	filename := fmt.Sprintf("%s.png", perkStyle)
	cachePath := filepath.Join(s.cacheDir, "perks", filename)

	// 检查缓存
	if data, err := s.readCache(cachePath); err == nil {
		return data, "image/png", nil
	}

	// 从CDN下载
	url := fmt.Sprintf("%s/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/%s/%s.png", CDRAGON_CDN, perkStyle, perkStyle)
	data, err := s.downloadAndCache(url, cachePath)
	if err != nil {
		return nil, "", err
	}

	return data, "image/png", nil
}

// readCache 读取缓存文件
func (s *AssetsService) readCache(path string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return os.ReadFile(path)
}

// downloadAndCache 下载并缓存
func (s *AssetsService) downloadAndCache(url, cachePath string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download: %s, status: %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	s.mu.Lock()
	defer s.mu.Unlock()
	os.WriteFile(cachePath, data, 0644)

	return data, nil
}

// GetCacheDir 获取缓存目录
func (s *AssetsService) GetCacheDir() string {
	return s.cacheDir
}
