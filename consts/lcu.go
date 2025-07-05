package consts

const (
	LolUxProcessName = "LeagueClientUx.exe"
)

// 游戏阶段枚举
type GameFlowPhase string

const (
	PhaseNone        GameFlowPhase = "None"
	PhaseLobby       GameFlowPhase = "Lobby"
	PhaseMatchmaking GameFlowPhase = "Matchmaking"
	PhaseChampSelect GameFlowPhase = "ChampSelect"
	PhaseInGame      GameFlowPhase = "InProgress"
	PhaseEndOfGame   GameFlowPhase = "EndOfGame"
)
