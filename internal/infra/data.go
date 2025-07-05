package infra

import (
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/infra/lcu"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewSQLite,
	NewLCUClient,
)

type Data struct {
	ConfData *conf.Data
	DB       *gorm.DB
	LCU      *lcu.Client
}

// NewData 构建统一资源聚合体
func NewData(c *conf.Data, db *gorm.DB, client *lcu.Client) (*Data, func(), error) {
	return &Data{
		ConfData: c,
		DB:       db,
		LCU:      client,
	}, func() {}, nil
}

// NewSQLite 初始化 SQLite 数据库
func NewSQLite(c *conf.Data, logger log.Logger) (*gorm.DB, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "infra/sqlite"))

	db, err := gorm.Open(sqlite.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	helper.Info("SQLite connected.")

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
			helper.Info("SQLite closed.")
		}
	}

	return db, cleanup, nil
}
func NewLCUClient(logger log.Logger, global *conf.Global) *lcu.Client {
	return lcu.NewClient(logger, global)
}
