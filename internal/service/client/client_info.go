package client

import (
	"encoding/json"
	"fmt"
	clientUseCase "github.com/B022MC/soraka-backend/internal/biz/client"
	"github.com/gin-gonic/gin"
	"go-utils/utils/ecode"
	"go-utils/utils/response"
	"net/http"
	"time"
)

type ClientInfoService struct {
	uc *clientUseCase.ClientInfoUseCase
}

func NewClientInfoService(uc *clientUseCase.ClientInfoUseCase) *ClientInfoService {
	return &ClientInfoService{uc: uc}
}

func (s *ClientInfoService) RegisterRouter(rootRouter *gin.RouterGroup) {
	privateRouter := rootRouter.Group("/client")
	privateRouter.GET("/getClientInfo", s.GetClientInfo)
	privateRouter.GET("/streamClientInfo", s.StreamClientInfo)

}

// GetClientInfo
// @Summary         获取客户端信息
// @Description     获取当前连接状态、游戏阶段、端口等信息
// @Tags            client/ClientInfo
// @Produce         json
// @Success         200 {object} response.Body{data=resp.ClientInfoResp,msg=string}
// @Failure         500 {object} response.Body{msg=string}
// @Router          /client/getClientInfo [get]
func (s *ClientInfoService) GetClientInfo(ctx *gin.Context) {
	clientInfo, err := s.uc.GetClientInfo()
	if err != nil {
		response.Fail(ctx, ecode.Failed, err.Error())
		return
	}
	response.Success(ctx, clientInfo)
}

// StreamClientInfo
// @Summary         实时推送客户端信息（SSE）
// @Description     每隔 2 秒推送一次客户端状态，基于 Server-Sent Events（EventStream）
// @Tags            client/ClientInfo
// @Produce         text/event-stream
// @Success         200 {string} string "data: {JSON}"
// @Failure         500 {string} string "Streaming unsupported"
// @Router          /client/streamClientInfo [get]
func (s *ClientInfoService) StreamClientInfo(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ticker := time.NewTicker(2 * time.Second)     // 推送数据
	heartbeat := time.NewTicker(10 * time.Second) // 发送心跳
	defer ticker.Stop()
	defer heartbeat.Stop()
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		http.Error(ctx.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	notify := ctx.Request.Context().Done()
	for {
		select {
		case <-notify:
			fmt.Printf("SSE 客户端断开连接")
			return
		case <-ticker.C:
			info, err := s.uc.GetClientInfo()
			if err != nil {
				_, writeErr := fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
				if writeErr != nil {
					fmt.Printf("SSE 写入失败: %v", writeErr)
					return
				}
				flusher.Flush()
				continue
			}
			jsonStr, _ := json.Marshal(info)
			_, writeErr := fmt.Fprintf(ctx.Writer, "data: %s\n\n", jsonStr)
			if writeErr != nil {
				fmt.Printf("SSE 写入失败: %v", writeErr)
				return
			}
			flusher.Flush()
		case <-heartbeat.C:
			_, err := fmt.Fprint(ctx.Writer, "event: ping\ndata: heartbeat\n\n")
			if err != nil {
				fmt.Printf("心跳写入失败: %v", err)
				return
			}
			flusher.Flush()
		}
	}
}
