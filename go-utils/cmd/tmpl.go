package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-utils/utils/lang/stringx"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	models          string
	to              string
	defaultTempPath string
	targetPkgPath   string
	folders         = []string{"dal/req", "dal/vo", "dal/repo", "biz", "service"}
	fileMap         = map[string]string{
		"dal/req":  "req",
		"dal/vo":   "vo",
		"dal/repo": "repo",
		"biz":      "biz",
		"service":  "service",
	}
)

func init() {
	TmplGeneratorCmd.Flags().StringVar(&models, "models", "", "eg: --models BaseUser")
	TmplGeneratorCmd.Flags().StringVar(&to, "to", "", "eg: --to biz/repo/req/vo/service")
	TmplGeneratorCmd.Flags().StringVar(&defaultTempPath, "defaultTempPath", "pkg/cmd/template", "eg: --defaultTempPath pkg/cmd/template")
	TmplGeneratorCmd.Flags().StringVar(&targetPkgPath, "targetPkgPath", "./internal", "eg: --targetPkgPath ../internal")
}

type TemplateData struct {
	ProjectName    string
	PkgName        string
	ModelName      string
	SnakeModelName string
}

var (
	TmplGeneratorCmd = &cobra.Command{
		Use:   "tmpl",
		Short: "tmpl",
		Long:  `tmpl.`,
		Run: func(cmd *cobra.Command, args []string) {
			if models == "" {
				panic("models is empty")
			}
			modelList := strings.Split(models, ",")
			pkgDir, err := os.Getwd()
			currentDir := filepath.Dir(pkgDir)
			if err != nil {
				panic(err)
			}
			genFolds := folders
			if to != "" {
				newFolds := make([]string, 0, len(folders))
				for _, folder := range folders {
					if strings.Contains(folder, to) {
						newFolds = append(newFolds, folder)
					}
				}
				genFolds = newFolds
			}

			for _, modelName := range modelList {
				for _, folder := range genFolds {
					templatePath := path.Join(currentDir, defaultTempPath, fmt.Sprintf("%s.go.tpl", fileMap[folder]))
					lowerCamel := stringx.FirstLower(modelName)
					snake := stringx.Camel2Snake(modelName)
					data := TemplateData{
						ProjectName:    project,
						PkgName:        pkg,
						ModelName:      modelName,
						SnakeModelName: lowerCamel,
					}
					// 读取模板文件
					tpl, err := template.ParseFiles(templatePath)
					if err != nil {
						panic(err)
					}

					modelPkgPath := path.Join(currentDir, targetPkgPath, folder, pkg)

					// 创建输出目录
					if _, err := os.Stat(modelPkgPath); os.IsNotExist(err) {
						os.Mkdir(modelPkgPath, 0755)
					}

					// 创建输出文件
					file, err := os.Create(fmt.Sprintf("%s/%s.go", modelPkgPath, snake))
					if err != nil {
						panic(err)
					}
					defer file.Close()

					// 执行模板并写入到文件
					if err := tpl.Execute(file, data); err != nil {
						panic(err)
					}
					fmt.Printf("生成文件成功: %s/%s.go\n", modelPkgPath, snake)
				}

			}

		},
	}
)
