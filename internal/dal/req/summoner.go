package req

type SummonerReq struct {
	Name  string `json:"name" form:"name"`
	Puuid string `json:"puuid" form:"puuid"`
}
