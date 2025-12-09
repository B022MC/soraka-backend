package resp

type GameDetail struct {
	EndOfGameResult  string `json:"endOfGameResult"`
	GameId           int64  `json:"gameId"`
	GameCreationDate string `json:"gameCreationDate"`
	GameDuration     int    `json:"gameDuration"`
	QueueId          int    `json:"queueId"`
	MapId            int    `json:"mapId"`
	GameMode         string `json:"gameMode"`
	GameType         string `json:"gameType"`
	Teams            []struct {
		TeamId          int    `json:"teamId"`
		Win             string `json:"win"`
		FirstBlood      bool   `json:"firstBlood"`
		FirstTower      bool   `json:"firstTower"`
		FirstInhibitor  bool   `json:"firstInhibitor"`
		FirstBaron      bool   `json:"firstBaron"`
		FirstDragon     bool   `json:"firstDragon"`
		FirstRiftHerald bool   `json:"firstRiftHerald"`
		TowerKills      int    `json:"towerKills"`
		InhibitorKills  int    `json:"inhibitorKills"`
		BaronKills      int    `json:"baronKills"`
		DragonKills     int    `json:"dragonKills"`
		RiftHeraldKills int    `json:"riftHeraldKills"`
		Bans            []struct {
			ChampionId int `json:"championId"`
			PickTurn   int `json:"pickTurn"`
		} `json:"bans"`
	} `json:"teams"`
	ParticipantIdentities []struct {
		Player struct {
			AccountId    int    `json:"accountId"`
			Puuid        string `json:"puuid"`
			PlatformId   string `json:"platformId"`
			SummonerName string `json:"summonerName"`
			GameName     string `json:"gameName"`
			TagLine      string `json:"tagLine"`
			SummonerId   int    `json:"summonerId"`
		} `json:"player"`
	} `json:"participantIdentities"`
	Participants []struct {
		ChampionKey   string `json:"championKey"`
		ParticipantId int    `json:"participantId"`
		TeamId        int    `json:"teamId"`
		ChampionId    int    `json:"championId"`
		Spell1Id      int    `json:"spell1Id"`
		Spell1Key     string `json:"spell1Key"`
		Spell2Id      int    `json:"spell2Id"`
		Spell2Key     string `json:"spell2Key"`
		Stats         struct {
			Win                 bool   `json:"win"`
			Item0               int    `json:"item0"`
			Item1               int    `json:"item1"`
			Item2               int    `json:"item2"`
			Item3               int    `json:"item3"`
			Item4               int    `json:"item4"`
			Item5               int    `json:"item5"`
			Item6               int    `json:"item6"`
			Item0Key            string `json:"item0Key"`
			Item1Key            string `json:"item1Key"`
			Item2Key            string `json:"item2Key"`
			Item3Key            string `json:"item3Key"`
			Item4Key            string `json:"item4Key"`
			Item5Key            string `json:"item5Key"`
			Item6Key            string `json:"item6Key"`
			Perk0               int    `json:"perk0"`
			Perk1               int    `json:"perk1"`
			Perk2               int    `json:"perk2"`
			Perk3               int    `json:"perk3"`
			Perk4               int    `json:"perk4"`
			Perk5               int    `json:"perk5"`
			PerkPrimaryStyle    int    `json:"perkPrimaryStyle"`
			PerkSubStyle        int    `json:"perkSubStyle"`
			PerkPrimaryStyleKey string `json:"perkPrimaryStyleKey"`
			PerkSubStyleKey     string `json:"perkSubStyleKey"`
			ChampLevel          int    `json:"champLevel"`

			Kills   int `json:"kills"`
			Deaths  int `json:"deaths"`
			Assists int `json:"assists"`

			GoldEarned                  int `json:"goldEarned"`
			GoldSpent                   int `json:"goldSpent"`
			TotalDamageDealtToChampions int `json:"totalDamageDealtToChampions"` //对英雄伤害
			TotalDamageDealt            int `json:"totalDamageDealt"`
			TotalDamageTaken            int `json:"totalDamageTaken"` //承受伤害
			TotalHeal                   int `json:"totalHeal"`
			TotalMinionsKilled          int `json:"totalMinionsKilled"`
		} `json:"stats"`
	} `json:"participants"`
}
