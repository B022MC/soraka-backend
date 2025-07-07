package resp

// QueueInfo 表示一个玩家在特定队列中的信息。
type QueueInfo struct {
	// QueueType 表示队列类型，例如 "RANKED_SOLO_5x5"。
	QueueType   string `json:"queueType"`
	QueueTypeCn string `json:"queueTypeCn"`

	// Division 表示玩家当前段位的分段，例如 "I"、"II"。
	Division string `json:"division"`
	Tier     string `json:"tier"`
	TierCn   string `json:"tierCn"`

	// HighestDivision 表示玩家历史最高的分段。
	HighestDivision string `json:"highestDivision"`

	// HighestTier 表示玩家历史最高的段位，例如 "Diamond"、"Master"。
	HighestTier string `json:"highestTier"`

	// IsProvisional 表示该队列是否处于定级赛阶段。
	IsProvisional bool `json:"isProvisional"`

	// LeaguePoints 表示玩家当前的段位点数（LP）。
	LeaguePoints int `json:"leaguePoints"`

	// Losses 表示玩家在该队列的失败场次。
	Losses int `json:"losses"`

	// Wins 表示玩家在该队列的胜利场次。
	Wins int `json:"wins"`
}

type QueueMap struct {
	RankedSolo5x5 QueueInfo `json:"RANKED_SOLO_5x5"`
	RankedFlexSr  QueueInfo `json:"RANKED_FLEX_SR"`
}
type Rank struct {
	QueueMap QueueMap `json:"queueMap"`
}
