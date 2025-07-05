package lcu

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) DoRequest(method, path string, body any) ([]byte, error) {
	c.mu.RLock()
	port := c.Port
	token := c.Token
	httpClient := c.httpClient
	c.mu.RUnlock()

	if port == 0 || token == "" {
		return nil, fmt.Errorf("LCU 客户端未连接")
	}

	url := fmt.Sprintf("https://127.0.0.1:%d%s", port, path)

	var reader io.Reader
	if body != nil {
		bts, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("请求参数序列化失败: %w", err)
		}
		reader = bytes.NewReader(bts)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	auth := base64.StdEncoding.EncodeToString([]byte("riot:" + token))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败（%s）: %w", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败 %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}
