package lcu

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/B022MC/soraka-backend/consts"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/local-utils/windows/process"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	lolCommandlineReg     = regexp.MustCompile(`--remoting-auth-token=(.+?)" ".*?--app-port=(\d+)"`)
	ErrLolProcessNotFound = errors.New("没有找到lol客户端进程")
)

type Client struct {
	mu         sync.RWMutex
	Connected  bool
	GamePhase  consts.GameFlowPhase
	SummonerID int64
	Token      string
	Port       int
	Host       string
	Polling    bool
	log        *log.Helper
	quit       chan struct{}
}

func NewClient(logger log.Logger, c *conf.Global) *Client {
	client := &Client{
		Host: c.Lcu.BaseUrl,
		log:  log.NewHelper(log.With(logger, "module", "infra/lcu")),
		quit: make(chan struct{}),
	}
	client.log.Info("LCU Client 初始化完成，开始轮询检测客户端进程")
	go client.backgroundLoop()
	return client
}

func (c *Client) backgroundLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.tryDetectClient()
		case <-c.quit:
			c.log.Info("停止 LCU 客户端检测")
			return
		}
	}
}

func (c *Client) tryDetectClient() {
	port, token, err := GetLolClientApiInfo()

	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		if c.Connected {
			c.log.Warn("LCU 客户端已关闭")
		}
		c.setDisconnectedLocked()
		return
	}

	if !c.Connected {
		c.log.Infof("检测到 LCU，端口: %d", port)
		go c.StartPolling()
	}

	c.Connected = true
	c.Port = port
	c.Token = token
}

func (c *Client) StartPolling() {
	c.mu.Lock()
	if c.Polling {
		c.mu.Unlock()
		return
	}
	c.Polling = true
	c.mu.Unlock()

	c.log.Info("启动游戏状态轮询")
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.pollStatus()
			case <-c.quit:
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

	resp, err := c.DoRequest(http.MethodGet, "/lol-gameflow/v1/gameflow-phase", nil)
	if err != nil {
		c.log.Errorf("拉取游戏阶段失败: %v", err)
		c.setDisconnected()
		return
	}

	var phase consts.GameFlowPhase
	if err := json.Unmarshal(resp, &phase); err != nil {
		c.log.Warnf("响应解析失败: %s", string(resp))
		return
	}

	c.mu.Lock()
	c.GamePhase = phase
	c.mu.Unlock()

	c.log.Infof("当前游戏阶段：%s", phase)
}

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
	c.GamePhase = consts.PhaseNone
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

func GetLolClientApiInfo() (port int, token string, err error) {
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

func (c *Client) DoRequest(method, path string, body any) ([]byte, error) {
	c.mu.RLock()
	port := c.Port
	token := c.Token
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

	cli := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
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
