package {{.PkgName}}

import (
	"github.com/go-kratos/kratos/v2/log"
	repox "go-utils/plugin/gormx/repo"
	{{.PkgName}}Model "{{.ProjectName}}/internal/dal/model/{{.PkgName}}"
	"{{.ProjectName}}/internal/infra"
)

type {{.ModelName}}Repo interface {
	repox.CORMImpl[{{.PkgName}}Model.{{.ModelName}}]
}

type {{.SnakeModelName}}Repo struct {
	repox.CORMImpl[{{.PkgName}}Model.{{.ModelName}}]
	data *infra.Data
	log  *log.Helper
}

// New{{.ModelName}}Repo .
func New{{.ModelName}}Repo(data *infra.Data, logger log.Logger) {{.ModelName}}Repo {

	return &{{.SnakeModelName}}Repo{
		CORMImpl: repox.NewCORMImplRepo[{{.PkgName}}Model.{{.ModelName}}](data),
		data:     data,
		log:      log.NewHelper(log.With(logger, "module", "repo/{{.SnakeModelName}}")),
	}
}

