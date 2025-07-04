package cmd

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go-utils/plugin/gormx/migrate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

var (
	migrateModels    string
	excludeModels    string
	dsn              string
	defaultModelPath string
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate",
		Long:  `migrate. 迁移`,
		Run: func(cmd *cobra.Command, args []string) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("failed to connect database: %v", err)
			}

			pkgDir, err := os.Getwd()
			currentDir := filepath.Dir(pkgDir)

			// 扫描models目录
			modelsDir := filepath.Join(currentDir, defaultModelPath, pkg)

			migrateModelList := strings.Split(migrateModels, ",")
			excludeModelList := strings.Split(excludeModels, ",")
			excludeMap := lo.SliceToMap(excludeModelList, func(item string) (string, struct{}) {
				return item, struct{}{}
			})
			migrateMap := lo.SliceToMap(migrateModelList, func(item string) (string, struct{}) {
				return item, struct{}{}
			})
			notFilter := false
			if len(migrateModelList) == 0 && len(excludeModelList) == 0 {
				notFilter = true
			}
			if err := migrate.ScanModels(modelsDir, func(name string) bool {
				if notFilter {
					return false
				}
				if _, ok := excludeMap[name]; ok {
					return true
				}
				if _, ok := migrateMap[name]; len(migrateMap) != 0 && ok {
					return false
				}
				return false
			}); err != nil {
				fmt.Sprintf("failed to scan models: %v", err)
			}

			// 执行自动迁移
			if err := migrate.AutoMigrate(db); err != nil {
				fmt.Sprintf("failed to migrate database: %v", err)
			}

		},
	}
)

func init() {
	migrateCmd.Flags().StringVar(&migrateModels, "migrateModels", "", "eg: --migrateModels BaseUser")
	migrateCmd.Flags().StringVar(&excludeModels, "excludeModels", "", "eg: --excludeModels BaseMenu")
	migrateCmd.Flags().StringVar(&dsn, "dsn", "", "eg: --dsn ")
	migrateCmd.Flags().StringVar(&defaultModelPath, "defaultModelPath", "internal/dal/model", "eg: --defaultModelPath internal/dal/model")

	rootCmd.AddCommand(migrateCmd)
}
