package champ_select

import (
	"github.com/B022MC/soraka-backend/internal/dal/repo/champ_select"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/go-kratos/kratos/v2/log"
)

type ChampSelectUseCase struct {
	repo champ_select.ChampSelectRepo
	log  *log.Helper
}

func NewChampSelectUseCase(repo champ_select.ChampSelectRepo, logger log.Logger) *ChampSelectUseCase {
	return &ChampSelectUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/champ_select")),
	}
}

func (uc *ChampSelectUseCase) GetSession() (*resp.ChampSelectSession, error) {
	return uc.repo.GetSession()
}

func (uc *ChampSelectUseCase) SelectChampion(actionId int64, championId int, completed bool) error {
	uc.log.Infof("Selecting champion %d for action %d", championId, actionId)
	return uc.repo.SelectChampion(actionId, championId, completed)
}

func (uc *ChampSelectUseCase) BanChampion(actionId int64, championId int, completed bool) error {
	uc.log.Infof("Banning champion %d for action %d", championId, actionId)
	return uc.repo.BanChampion(actionId, championId, completed)
}

func (uc *ChampSelectUseCase) AcceptTrade(tradeId int64) error {
	return uc.repo.AcceptTrade(tradeId)
}

func (uc *ChampSelectUseCase) AcceptSwap(swapId int64) error {
	return uc.repo.AcceptSwap(swapId)
}

func (uc *ChampSelectUseCase) BenchSwap(championId int) error {
	return uc.repo.BenchSwap(championId)
}

func (uc *ChampSelectUseCase) GetCurrentChampion() (int, error) {
	return uc.repo.GetCurrentChampion()
}

func (uc *ChampSelectUseCase) GetSkinCarousel() ([]resp.SkinCarousel, error) {
	return uc.repo.GetSkinCarousel()
}

func (uc *ChampSelectUseCase) SelectSkin(skinId int, spell1Id, spell2Id *int) error {
	return uc.repo.SelectSkin(skinId, spell1Id, spell2Id)
}

func (uc *ChampSelectUseCase) Reroll() error {
	return uc.repo.Reroll()
}
