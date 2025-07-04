package {{.PkgName}}

import {{.PkgName}}Model "{{.ProjectName}}/internal/dal/model/{{.PkgName}}"

type {{.ModelName}}Vo struct {
	*{{.PkgName}}Model.{{.ModelName}} `mapstructure:",squash"`
	// eg: Name string `json:"name" from:"name" mapstructure:"name" title:"名称" binding:"required"`
}
