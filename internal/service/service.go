package service

import (
	LcuService "github.com/B022MC/soraka-backend/internal/service/lcu"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	LcuService.NewCurrentSummonerService,
	LcuService.NewGamePhaseService,
)
