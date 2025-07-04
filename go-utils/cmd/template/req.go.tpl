package {{.PkgName}}

import "go-utils/utils"

type Add{{.ModelName}}Req struct {
	// eg: Name string `json:"name" from:"name" mapstructure:"name" title:"名称" binding:"required"`

}

type Update{{.ModelName}}Req struct {
	Add{{.ModelName}}Req `mapstructure:",squash"`
	Id int32 `json:"id" from:"id" mapstructure:"id" title:"主键" binding:"required"` // 主键
}

type List{{.ModelName}}Req struct {
	*utils.PageParam
}
