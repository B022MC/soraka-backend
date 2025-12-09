package resp

// GameflowSession 游戏流程会话信息
type GameflowSession struct {
	Phase    string `json:"phase"`
	GameData struct {
		GameId  int64        `json:"gameId"`
		Queue   Queue        `json:"queue"`
		TeamOne []TeamMember `json:"teamOne"`
		TeamTwo []TeamMember `json:"teamTwo"`
	} `json:"gameData"`
}

// TeamMember 队伍成员信息
type TeamMember struct {
	SummonerId   int64  `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	Puuid        string `json:"puuid"`
	ChampionId   int    `json:"championId"`
}

// Queue 队列信息
type Queue struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ReadyCheckStatus 准备确认状态
type ReadyCheckStatus struct {
	State          string  `json:"state"`
	PlayerResponse string  `json:"playerResponse"`
	DodgeWarning   string  `json:"dodgeWarning"`
	Timer          float64 `json:"timer"`
	DeclinerIds    []int64 `json:"declinerIds"`
}
