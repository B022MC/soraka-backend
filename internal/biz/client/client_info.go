package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	clientRepo "github.com/B022MC/soraka-backend/internal/dal/repo/client"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ClientInfoUseCase struct {
	repo clientRepo.ClientInfoRepo
	log  *log.Helper
}

func NewClientUseCase(repo clientRepo.ClientInfoRepo, logger log.Logger) *ClientInfoUseCase {
	return &ClientInfoUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "uc/client_info")),
	}
}

func (c *ClientInfoUseCase) GetClientInfo() (*resp.ClientInfoResp, error) {
	result, err := c.repo.GetClientInfo()
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	return result, nil
}
func (c *ClientInfoUseCase) OpenLolClient() error {
	return c.repo.OpenLolClient()
}

func (c *ClientInfoUseCase) ReverseProxy(ctx *gin.Context) error {
	clientInfo, err := c.repo.GetClientInfo()
	if err != nil {
		return err
	}
	if !clientInfo.Connected {
		return errors.New("LCU 未连接")
	}

	rawURL := fmt.Sprintf("https://127.0.0.1:%d", clientInfo.Port)
	targetURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("LCU 地址解析失败: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 重写 Director 设置 Basic Auth
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.RequestURI = ""
		// 去掉前缀 /client/proxy
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/client/proxy")
		req.SetBasicAuth("riot", clientInfo.Token)
	}

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
	return nil
}
