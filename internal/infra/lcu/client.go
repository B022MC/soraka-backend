// client.go
package lcu

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
)

type Client struct {
	mu         sync.RWMutex
	Connected  bool
	GamePhase  string
	Token      string
	Port       int
	Host       string
	Polling    bool
	log        *log.Helper
	quit       chan struct{}
	failCount  int
	httpClient *http.Client
	initOnce   sync.Once
	conf       *conf.Global
	ClientPath string
}

func NewClient(logger log.Logger, conf *conf.Global) *Client {
	client := &Client{
		Host:       conf.Lcu.BaseUrl,
		log:        log.NewHelper(log.With(logger, "module", "infra/lcu")),
		quit:       make(chan struct{}),
		httpClient: newHttpClient(),
		conf:       conf,
	}

	client.log.Info("LCU Client 初始化完成，开始检测客户端进程")
	go client.backgroundLoop()
	return client
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
