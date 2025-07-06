package service

import (
	clientService "github.com/B022MC/soraka-backend/internal/service/client"
	lcuService "github.com/B022MC/soraka-backend/internal/service/lcu"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	lcuService.NewCurrentSummonerService,
	clientService.NewClientInfoService,
)
