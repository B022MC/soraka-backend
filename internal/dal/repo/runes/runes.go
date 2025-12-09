package runes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type RunesRepo interface {
	GetCurrentPage() (*resp.RunePage, error)
	DeletePage(pageId int) error
	CreatePage(name string, primaryStyleId, subStyleId int, selectedPerkIds []int) error
	GetAllPages() ([]resp.RunePage, error)
}

type runesRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewRunesRepo(data *infra.Data, global *conf.Global, logger log.Logger) RunesRepo {
	return &runesRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/runes")),
	}
}

func (r *runesRepo) GetCurrentPage() (*resp.RunePage, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.PerksCurrentPagePath, nil)
	if err != nil {
		return nil, err
	}

	var page resp.RunePage
	if err := json.Unmarshal(request, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

func (r *runesRepo) DeletePage(pageId int) error {
	uri := fmt.Sprintf("%s/%d", r.global.Lcu.PerksPagesPath, pageId)
	_, err := r.data.LCU.DoRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to delete rune page: %w", err)
	}

	r.log.Infof("Rune page %d deleted", pageId)
	return nil
}

func (r *runesRepo) CreatePage(name string, primaryStyleId, subStyleId int, selectedPerkIds []int) error {
	data := map[string]interface{}{
		"name":            name,
		"primaryStyleId":  primaryStyleId,
		"subStyleId":      subStyleId,
		"selectedPerkIds": selectedPerkIds,
		"current":         true,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.PerksPagesPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create rune page: %w", err)
	}

	r.log.Infof("Rune page '%s' created", name)
	return nil
}

func (r *runesRepo) GetAllPages() ([]resp.RunePage, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.PerksPagesPath, nil)
	if err != nil {
		return nil, err
	}

	var pages []resp.RunePage
	if err := json.Unmarshal(request, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}
