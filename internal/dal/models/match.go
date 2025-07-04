package models

type MatchRecord struct {
	ID         int32  `json:"id"`
	GameID     string `json:"game_id"`
	SummonerID string
	Win        bool
	CreatedAt  int64
}
