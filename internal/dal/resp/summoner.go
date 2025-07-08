package resp

import (
	"github.com/B022MC/soraka-backend/consts"
)

type Summoner struct {
	GameName                    string  `json:"gameName"`
	TagLine                     string  `json:"tagLine"`
	SummonerLevel               int     `json:"summonerLevel"`
	ProfileIconId               int     `json:"profileIconId"`
	ProfileIconUrl              string  `json:"profileIconUrl"`
	Puuid                       string  `json:"puuid"`
	PlatformIdCn                string  `json:"platformIdCn"`
	XpSinceLastLevel            int     `json:"xpSinceLastLevel"`
	XpUntilNextLevel            int     `json:"xpUntilNextLevel"`
	PercentCompleteForNextLevel float64 `json:"percentCompleteForNextLevel"`
}
type SummonerInfo struct {
	Summoner Summoner `json:"summoner"`
	Rank     Rank     `json:"rank"`
}

func (s *Summoner) Dto() {
	// 根据 profile icon ID 构造 URL（从 consts 中获取图标路径）
	if path, ok := consts.ProfileIconMap[s.ProfileIconId]; ok {
		s.ProfileIconUrl = path
	} else {
		// fallback: 默认图标或空
		s.ProfileIconUrl = ""
	}
}
