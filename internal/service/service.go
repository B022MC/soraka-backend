package service

import (
	"github.com/google/wire"
	LcuService "soraka-backend/internal/service/lcu"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	LcuService.NewCurrentSummonerService,
)
