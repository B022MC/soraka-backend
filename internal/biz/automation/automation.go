package automation

import (
	"github.com/B022MC/soraka-backend/internal/dal/repo/champ_select"
	"github.com/B022MC/soraka-backend/internal/dal/repo/gameflow"
	"github.com/B022MC/soraka-backend/internal/dal/repo/profile"
	"github.com/B022MC/soraka-backend/internal/dal/repo/runes"
	"github.com/go-kratos/kratos/v2/log"
)

// AutomationUseCase 自动化操作用例
// 用于实现自动BP、自动接受对局等功能
type AutomationUseCase struct {
	gameflowRepo    gameflow.GameflowRepo
	champSelectRepo champ_select.ChampSelectRepo
	runesRepo       runes.RunesRepo
	profileRepo     profile.ProfileRepo
	log             *log.Helper
}

func NewAutomationUseCase(
	gameflowRepo gameflow.GameflowRepo,
	champSelectRepo champ_select.ChampSelectRepo,
	runesRepo runes.RunesRepo,
	profileRepo profile.ProfileRepo,
	logger log.Logger,
) *AutomationUseCase {
	return &AutomationUseCase{
		gameflowRepo:    gameflowRepo,
		champSelectRepo: champSelectRepo,
		runesRepo:       runesRepo,
		profileRepo:     profileRepo,
		log:             log.NewHelper(log.With(logger, "module", "biz/automation")),
	}
}

// AutoAcceptReadyCheck 自动接受对局
func (uc *AutomationUseCase) AutoAcceptReadyCheck() error {
	status, err := uc.gameflowRepo.GetReadyCheckStatus()
	if err != nil {
		return err
	}

	if status.State == "InProgress" && status.PlayerResponse == "None" {
		return uc.gameflowRepo.AcceptReadyCheck()
	}

	return nil
}

// AutoSelectChampion 自动选择英雄
func (uc *AutomationUseCase) AutoSelectChampion(championId int) error {
	session, err := uc.champSelectRepo.GetSession()
	if err != nil {
		return err
	}

	// 查找当前玩家的选择动作
	for _, actionGroup := range session.Actions {
		for _, action := range actionGroup {
			if action.ActorCellId == session.LocalPlayerCellId &&
				action.Type == "pick" &&
				!action.Completed &&
				action.IsInProgress {
				return uc.champSelectRepo.SelectChampion(action.Id, championId, true)
			}
		}
	}

	return nil
}

// AutoBanChampion 自动禁用英雄
func (uc *AutomationUseCase) AutoBanChampion(championId int) error {
	session, err := uc.champSelectRepo.GetSession()
	if err != nil {
		return err
	}

	// 查找当前玩家的禁用动作
	for _, actionGroup := range session.Actions {
		for _, action := range actionGroup {
			if action.ActorCellId == session.LocalPlayerCellId &&
				action.Type == "ban" &&
				!action.Completed &&
				action.IsInProgress {
				return uc.champSelectRepo.BanChampion(action.Id, championId, true)
			}
		}
	}

	return nil
}

// AutoAcceptTrades 自动接受英雄交换
func (uc *AutomationUseCase) AutoAcceptTrades() error {
	session, err := uc.champSelectRepo.GetSession()
	if err != nil {
		return err
	}

	for _, trade := range session.Trades {
		if trade.State == "RECEIVED" {
			if err := uc.champSelectRepo.AcceptTrade(trade.Id); err != nil {
				uc.log.Errorf("Failed to accept trade %d: %v", trade.Id, err)
			}
		}
	}

	return nil
}

// AutoAcceptSwaps 自动接受楼层交换
func (uc *AutomationUseCase) AutoAcceptSwaps() error {
	session, err := uc.champSelectRepo.GetSession()
	if err != nil {
		return err
	}

	for _, swap := range session.PickOrderSwaps {
		if swap.State == "RECEIVED" {
			if err := uc.champSelectRepo.AcceptSwap(swap.Id); err != nil {
				uc.log.Errorf("Failed to accept swap %d: %v", swap.Id, err)
			}
		}
	}

	return nil
}

// ApplyRunePage 应用符文页（用于OPGG符文一键设置）
func (uc *AutomationUseCase) ApplyRunePage(name string, primaryStyleId, subStyleId int, selectedPerkIds []int) error {
	// 获取当前符文页
	currentPage, err := uc.runesRepo.GetCurrentPage()
	if err == nil && currentPage.IsDeletable {
		// 删除当前可删除的符文页
		_ = uc.runesRepo.DeletePage(currentPage.Id)
	}

	// 创建新符文页
	return uc.runesRepo.CreatePage(name, primaryStyleId, subStyleId, selectedPerkIds)
}
