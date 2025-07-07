package req

type MatchHistoryReq struct {
	Puuid    string `json:"puuid" form:"puuid"`
	BegIndex int    `json:"beg_index" form:"beg_index"`
	EndIndex int    `json:"end_index" form:"end_index"`
}
