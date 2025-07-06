package resp

type ClientInfoResp struct {
	Connected  bool   `json:"connected" example:"true"`                // 是否连接
	GamePhase  string `json:"game_phase" example:"ChampSelect"`        // 当前游戏阶段
	Token      string `json:"token" example:"xyz123"`                  // LCU Token
	Port       int    `json:"port" example:"12345"`                    // LCU 端口
	ClientPath string `json:"client_path" example:"D:/LOL/Client.exe"` // 客户端路径
}
