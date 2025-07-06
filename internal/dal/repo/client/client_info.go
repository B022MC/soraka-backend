package client

import (
	"errors"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type ClientInfoRepo interface {
	GetClientInfo() (*resp.ClientInfoResp, error)
}
type clientInfoRepo struct {
	data *infra.Data
	log  *log.Helper
}

func (c clientInfoRepo) GetClientInfo() (*resp.ClientInfoResp, error) {
	lcuClient := c.data.LCU
	if lcuClient == nil {
		c.log.Warn("LCU Client 实例未初始化")
		return nil, errors.New("LCU Client 实例未初始化")
	}

	lcuClient.Mutex().RLock()
	defer lcuClient.Mutex().RUnlock()

	return &resp.ClientInfoResp{
		Connected:  lcuClient.IsConnected(),
		GamePhase:  lcuClient.GetGamePhase(),
		Token:      lcuClient.GetToken(),
		Port:       lcuClient.GetPort(),
		ClientPath: lcuClient.GetClientPath(),
	}, nil
}

func NewClientInfoRepo(data *infra.Data, logger log.Logger) ClientInfoRepo {
	return &clientInfoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/client_info")),
	}
}
