package client

import (
	clientRepo "github.com/B022MC/soraka-backend/internal/dal/repo/client"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/go-kratos/kratos/v2/log"
)

type ClientInfoUseCase struct {
	repo clientRepo.ClientInfoRepo
	log  *log.Helper
}

func NewClientUseCase(repo clientRepo.ClientInfoRepo, logger log.Logger) *ClientInfoUseCase {
	return &ClientInfoUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "uc/client_info")),
	}
}

func (c *ClientInfoUseCase) GetClientInfo() (*resp.ClientInfoResp, error) {
	result, err := c.repo.GetClientInfo()
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	return result, nil
}
