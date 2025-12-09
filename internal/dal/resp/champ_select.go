package resp

// ChampSelectSession 英雄选择会话信息
type ChampSelectSession struct {
	Timer struct {
		AdjustedTimeLeftInPhase int64  `json:"adjustedTimeLeftInPhase"`
		InternalNowInEpochMs    int64  `json:"internalNowInEpochMs"`
		IsInfinite              bool   `json:"isInfinite"`
		Phase                   string `json:"phase"`
		TotalTimeInPhase        int64  `json:"totalTimeInPhase"`
	} `json:"timer"`
	ChatDetails struct {
		ChatRoomName     string `json:"chatRoomName"`
		ChatRoomPassword string `json:"chatRoomPassword"`
	} `json:"chatDetails"`
	MyTeam            []ChampSelectPlayer `json:"myTeam"`
	TheirTeam         []ChampSelectPlayer `json:"theirTeam"`
	Trades            []Trade             `json:"trades"`
	PickOrderSwaps    []Swap              `json:"pickOrderSwaps"`
	Actions           [][]Action          `json:"actions"`
	BenchChampionIds  []int               `json:"benchChampionIds"`
	LocalPlayerCellId int64               `json:"localPlayerCellId"`
}

// ChampSelectPlayer 英雄选择中的玩家信息
type ChampSelectPlayer struct {
	AssignedPosition   string `json:"assignedPosition"`
	CellId             int64  `json:"cellId"`
	ChampionId         int    `json:"championId"`
	ChampionPickIntent int    `json:"championPickIntent"`
	SelectedSkinId     int    `json:"selectedSkinId"`
	Spell1Id           int    `json:"spell1Id"`
	Spell2Id           int    `json:"spell2Id"`
	SummonerId         int64  `json:"summonerId"`
	Team               int    `json:"team"`
	WardSkinId         int    `json:"wardSkinId"`
}

// Action 选择/禁用动作
type Action struct {
	Id           int64  `json:"id"`
	ActorCellId  int64  `json:"actorCellId"`
	ChampionId   int    `json:"championId"`
	Type         string `json:"type"` // pick, ban
	Completed    bool   `json:"completed"`
	IsInProgress bool   `json:"isInProgress"`
}

// Trade 英雄交换请求
type Trade struct {
	Id     int64  `json:"id"`
	CellId int64  `json:"cellId"`
	State  string `json:"state"` // AVAILABLE, BUSY, INVALID, RECEIVED, SENT
}

// Swap 楼层交换请求
type Swap struct {
	Id     int64  `json:"id"`
	CellId int64  `json:"cellId"`
	State  string `json:"state"` // AVAILABLE, BUSY, INVALID, RECEIVED, SENT
}

// SkinCarousel 皮肤轮播
type SkinCarousel struct {
	ChampionId int    `json:"championId"`
	SkinId     int    `json:"id"`
	Name       string `json:"name"`
	Ownership  struct {
		Owned bool `json:"owned"`
	} `json:"ownership"`
	SplashPath           string `json:"splashPath"`
	UncenteredSplashPath string `json:"uncenteredSplashPath"`
	TilePath             string `json:"tilePath"`
}
