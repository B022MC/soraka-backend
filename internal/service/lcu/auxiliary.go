package lcu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"go-utils/utils/ecode"
	"go-utils/utils/response"

	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/windows/registry"
)

type AuxiliaryService struct {
	data *infra.Data
}

func NewAuxiliaryService(data *infra.Data) *AuxiliaryService {
	return &AuxiliaryService{data: data}
}

// 保存的客户端路径列表
var clientPaths = []string{}

func (s *AuxiliaryService) RegisterRouter(rootRouter *gin.RouterGroup) {
	r := rootRouter.Group("/lcu")
	r.GET("/clientPaths", s.GetClientPaths)
	r.GET("/customGameConfigs", s.GetCustomGameConfigs)
	r.POST("/startClient", s.StartClient)
	r.POST("/createPracticeLobby", s.CreatePracticeLobby)
	r.POST("/spectate", s.Spectate)
	r.POST("/restartClient", s.RestartClient)
	r.POST("/fixClientWindow", s.FixClientWindow)
	r.POST("/setOnlineStatus", s.SetOnlineStatus)
	r.POST("/setAvailability", s.SetAvailability)
	r.POST("/setProfileTier", s.SetProfileTier)
	r.POST("/setProfileBackground", s.SetProfileBackground)
	r.POST("/removeTokens", s.RemoveTokens)
	r.POST("/removePrestigeCrest", s.RemovePrestigeCrest)
}

// getLoLPathFromRegistry 从注册表获取LOL安装路径（国服）
func getLoLPathFromRegistry() string {
	// 国服注册表路径: HKEY_CURRENT_USER\SOFTWARE\Tencent\LOL
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Tencent\LOL`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	installPath, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return ""
	}

	// 标准化路径
	path := strings.ReplaceAll(installPath, "\\", "/")
	path = strings.TrimSuffix(path, "/")
	path = path + "/TCLS"
	return path
}

// GetClientPaths 获取已知的客户端路径列表
func (s *AuxiliaryService) GetClientPaths(ctx *gin.Context) {
	// 1. 从注册表获取路径
	regPath := getLoLPathFromRegistry()
	if regPath != "" {
		addClientPath(regPath)
	}

	// 2. 如果LCU已连接，尝试获取安装目录并添加到列表
	if s.data != nil && s.data.LCU != nil {
		res, err := s.data.LCU.DoRequest(http.MethodGet, "/data-store/v1/install-dir", nil)
		if err == nil && len(res) > 0 {
			var installDir string
			if err := json.Unmarshal(res, &installDir); err == nil && installDir != "" {
				// 标准化路径
				installDir = strings.ReplaceAll(installDir, "\\", "/")
				installDir = strings.ReplaceAll(installDir, "/LeagueClient", "/TCLS")
				addClientPath(installDir)
			}
		}
	}

	response.Success(ctx, clientPaths)
}

// AddClientPath 添加客户端路径到列表
func addClientPath(path string) {
	path = strings.ReplaceAll(path, "\\", "/")
	for _, p := range clientPaths {
		if strings.EqualFold(p, path) {
			return
		}
	}
	clientPaths = append(clientPaths, path)
}

// StartClient 启动LOL客户端
func (s *AuxiliaryService) StartClient(ctx *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	_ = ctx.ShouldBindJSON(&req)

	// 确保从注册表获取路径
	regPath := getLoLPathFromRegistry()
	if regPath != "" {
		addClientPath(regPath)
		fmt.Printf("从注册表获取到路径: %s\n", regPath)
	} else {
		fmt.Println("未能从注册表获取LOL路径")
	}
	fmt.Printf("当前已知路径列表: %v\n", clientPaths)

	// 如果指定了路径，使用指定路径
	if req.Path != "" {
		for _, clientName := range []string{"client.exe", "LeagueClient.exe"} {
			fullPath := req.Path + "/" + clientName
			if _, err := os.Stat(fullPath); err == nil {
				cmd := exec.Command(fullPath)
				if err := cmd.Start(); err == nil {
					addClientPath(req.Path)
					response.Success(ctx, map[string]string{"path": fullPath})
					return
				}
			}
		}
		response.Fail(ctx, ecode.Failed, "指定路径无效")
		return
	}

	// 没有指定路径，尝试从已知路径列表启动
	for _, basePath := range clientPaths {
		for _, clientName := range []string{"client.exe", "LeagueClient.exe"} {
			fullPath := basePath + "/" + clientName
			if _, err := os.Stat(fullPath); err == nil {
				cmd := exec.Command(fullPath)
				if err := cmd.Start(); err == nil {
					response.Success(ctx, map[string]string{"path": fullPath})
					return
				}
			}
		}
	}

	response.Fail(ctx, ecode.Failed, "未找到LOL客户端，请先连接一次客户端或手动启动")
}

// checkLCU 检查LCU客户端是否已连接
func (s *AuxiliaryService) checkLCU(ctx *gin.Context) bool {
	if s.data == nil || s.data.LCU == nil {
		response.Fail(ctx, ecode.Failed, "游戏客户端未连接")
		return false
	}
	return true
}

// GetCustomGameConfigs 获取可用的自定义游戏配置（调试用）
func (s *AuxiliaryService) GetCustomGameConfigs(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}

	// 获取可用的游戏模式
	res, err := s.data.LCU.DoRequest(http.MethodGet, "/lol-game-queues/v1/custom", nil)
	if err != nil {
		response.Fail(ctx, ecode.Failed, fmt.Sprintf("获取自定义游戏配置失败: %v", err))
		return
	}

	var configs interface{}
	json.Unmarshal(res, &configs)
	response.Success(ctx, configs)
}

// CreatePracticeLobby 创建5v5自定义房间
func (s *AuxiliaryService) CreatePracticeLobby(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Mode     string `json:"mode"` // PRACTICETOOL 或 CLASSIC
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Name == "" {
		response.Fail(ctx, ecode.CaptchaFailed, "房间名称不能为空")
		return
	}

	// 先退出当前房间/队列
	s.data.LCU.DoRequest(http.MethodDelete, "/lol-lobby/v2/lobby", nil)
	s.data.LCU.DoRequest(http.MethodDelete, "/lol-lobby/v2/lobby/matchmaking/search", nil)

	// 与Seraphine完全相同的配置
	body := map[string]interface{}{
		"customGameLobby": map[string]interface{}{
			"configuration": map[string]interface{}{
				"gameMode":         "PRACTICETOOL",
				"gameMutator":      "",
				"gameServerRegion": "",
				"mapId":            11,
				"mutators":         map[string]interface{}{"id": 1},
				"spectatorPolicy":  "AllAllowed",
				"teamSize":         5,
			},
			"lobbyName":     req.Name,
			"lobbyPassword": req.Password,
		},
		"isCustom": true,
	}

	_, err := s.data.LCU.DoRequest(http.MethodPost, "/lol-lobby/v2/lobby", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, fmt.Sprintf("创建练习房失败: %v (提示：国服可能不支持此功能，或账号等级不足)", err))
		return
	}
	response.Success(ctx, nil)
}

// Spectate 观战
func (s *AuxiliaryService) Spectate(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Name == "" {
		response.Fail(ctx, ecode.CaptchaFailed, "召唤师名称不能为空")
		return
	}

	// 1. 根据名称获取召唤师信息（使用params方式传参，自动URL编码）
	params := map[string]string{"name": req.Name}
	summonerRes, err := s.data.LCU.DoRequestWithParams(http.MethodGet, "/lol-summoner/v1/summoners", params, nil)
	if err != nil {
		response.Fail(ctx, ecode.Failed, fmt.Sprintf("获取召唤师信息失败: %v", err))
		return
	}

	var summoner struct {
		Puuid string `json:"puuid"`
	}
	if err := json.Unmarshal(summonerRes, &summoner); err != nil || summoner.Puuid == "" {
		response.Fail(ctx, ecode.Failed, "未找到该召唤师")
		return
	}

	// 2. 发起观战请求
	spectateData := map[string]interface{}{
		"allowObserveMode":     "ALL",
		"dropInSpectateGameId": req.Name,
		"gameQueueType":        "",
		"puuid":                summoner.Puuid,
	}
	res, err := s.data.LCU.DoRequest(http.MethodPost, "/lol-spectator/v1/spectate/launch", spectateData)
	if err != nil {
		response.Fail(ctx, ecode.Failed, fmt.Sprintf("观战请求失败: %v", err))
		return
	}
	if len(res) > 0 {
		response.Fail(ctx, ecode.Failed, "该召唤师当前不在游戏中")
		return
	}
	response.Success(ctx, nil)
}

// RestartClient 重启客户端
func (s *AuxiliaryService) RestartClient(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	_, err := s.data.LCU.DoRequest(http.MethodPost, "/riotclient/kill-and-restart-ux", nil)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "重启客户端失败")
		return
	}
	response.Success(ctx, nil)
}

// FixClientWindow 修复客户端窗口
func (s *AuxiliaryService) FixClientWindow(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	// 获取当前缩放值
	res, err := s.data.LCU.DoRequest(http.MethodGet, "/riotclient/zoom-scale", nil)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取缩放值失败")
		return
	}

	var zoom float64
	if err := json.Unmarshal(res, &zoom); err != nil {
		response.Fail(ctx, ecode.Failed, "解析缩放值失败")
		return
	}

	// 计算目标窗口大小
	var width, height int
	switch {
	case zoom <= 0.8:
		width, height = 1024, 576
	case zoom <= 1.0:
		width, height = 1280, 720
	default:
		width, height = 1600, 900
	}

	// 尝试调用Windows API修复窗口（需要管理员权限）
	// 通过调用riotclient接口来尝试修复
	fixData := map[string]interface{}{
		"width":  width,
		"height": height,
	}
	_, _ = s.data.LCU.DoRequest(http.MethodPost, "/riotclient/ux-minimize", nil)
	_, _ = s.data.LCU.DoRequest(http.MethodPost, "/riotclient/ux-show", nil)

	response.Success(ctx, map[string]interface{}{
		"zoom":   zoom,
		"width":  width,
		"height": height,
		"fixed":  fixData,
		"msg":    "已尝试刷新客户端窗口，如未生效请以管理员身份运行",
	})
}

// SetOnlineStatus 设置在线状态
func (s *AuxiliaryService) SetOnlineStatus(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	body := map[string]string{"statusMessage": req.Status}
	_, err := s.data.LCU.DoRequest(http.MethodPut, "/lol-chat/v1/me", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "设置在线状态失败")
		return
	}
	response.Success(ctx, nil)
}

// SetAvailability 设置在线可用性
func (s *AuxiliaryService) SetAvailability(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		Availability string `json:"availability"` // chat, away, offline
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	body := map[string]string{"availability": req.Availability}
	_, err := s.data.LCU.DoRequest(http.MethodPut, "/lol-chat/v1/me", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "设置可用性失败")
		return
	}
	response.Success(ctx, nil)
}

// SetProfileTier 设置段位展示
func (s *AuxiliaryService) SetProfileTier(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		Queue    string `json:"queue"`    // RANKED_SOLO_5x5, RANKED_FLEX_SR, RANKED_TFT
		Tier     string `json:"tier"`     // IRON, BRONZE, etc.
		Division string `json:"division"` // I, II, III, IV
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	body := map[string]interface{}{
		"lol": map[string]string{
			"rankedLeagueQueue":    req.Queue,
			"rankedLeagueTier":     req.Tier,
			"rankedLeagueDivision": req.Division,
		},
	}
	_, err := s.data.LCU.DoRequest(http.MethodPut, "/lol-chat/v1/me", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "设置段位展示失败")
		return
	}
	response.Success(ctx, nil)
}

// SetProfileBackground 设置个人背景皮肤
func (s *AuxiliaryService) SetProfileBackground(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	var req struct {
		SkinId int `json:"skinId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.SkinId == 0 {
		response.Fail(ctx, ecode.CaptchaFailed, "皮肤ID不能为空")
		return
	}

	body := map[string]interface{}{
		"key":   "backgroundSkinId",
		"value": req.SkinId,
	}
	_, err := s.data.LCU.DoRequest(http.MethodPost, "/lol-summoner/v1/current-summoner/summoner-profile", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "设置个人背景失败")
		return
	}
	response.Success(ctx, nil)
}

// RemoveTokens 移除挑战徽章
func (s *AuxiliaryService) RemoveTokens(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	body := map[string]interface{}{
		"challengeIds": []int{},
	}
	_, err := s.data.LCU.DoRequest(http.MethodPost, "/lol-challenges/v1/update-player-preferences", body)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "移除挑战徽章失败")
		return
	}
	response.Success(ctx, nil)
}

// RemovePrestigeCrest 移除声望边框
func (s *AuxiliaryService) RemovePrestigeCrest(ctx *gin.Context) {
	if !s.checkLCU(ctx) {
		return
	}
	// 获取当前设置
	res, err := s.data.LCU.DoRequest(http.MethodGet, "/lol-regalia/v2/current-summoner/regalia", nil)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取当前设置失败")
		return
	}

	var regalia map[string]interface{}
	if err := json.Unmarshal(res, &regalia); err != nil {
		response.Fail(ctx, ecode.Failed, "解析设置失败")
		return
	}

	// 设置 crestType 为 "prestige"，prestigeCrestLevel 为 1
	regalia["crestType"] = "prestige"
	regalia["selectedPrestigeCrest"] = 1

	_, err = s.data.LCU.DoRequest(http.MethodPut, "/lol-regalia/v2/current-summoner/regalia", regalia)
	if err != nil {
		response.Fail(ctx, ecode.Failed, fmt.Sprintf("移除声望边框失败: %v", err))
		return
	}
	response.Success(ctx, nil)
}
