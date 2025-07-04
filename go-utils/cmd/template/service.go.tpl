package {{.PkgName}}

import (
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"go-utils/plugin/middleware"
	"go-utils/utils"
	"go-utils/utils/ecode"
	"go-utils/utils/response"
	{{.PkgName}}Biz "{{.ProjectName}}/internal/biz/{{.PkgName}}"
	{{.PkgName}}Req "{{.ProjectName}}/internal/dal/req/{{.PkgName}}"
)

// {{.ModelName}}Service is a {{.ModelName}} service.
type {{.ModelName}}Service struct {
	uc *{{.PkgName}}Biz.{{.ModelName}}UseCase
}

// New{{.ModelName}}Service new a {{.ModelName}} service.
func New{{.ModelName}}Service(uc *{{.PkgName}}Biz.{{.ModelName}}UseCase) *{{.ModelName}}Service {
	return &{{.ModelName}}Service{uc: uc}
}

// @Description: RegisterRouter register router.注冊路由服务
// @receiver s
// @param publicRouter
// @param privateRouter
func (s *{{.ModelName}}Service) RegisterRouter(rootRouter *gin.RouterGroup) {

	//publicRouter := rootRouter.Group("/{{.PkgName}}/user").Use(middleware.SwitchingDB())
	privateRouter := rootRouter.Group("/{{.PkgName}}/user").Use(middleware.JWTAuth())
	privateRouter.POST("/addOne", s.AddOne)
	privateRouter.GET("/delOne", s.DelOne)
	privateRouter.POST("/delMany", s.DelMany)
	privateRouter.GET("/getOne", s.GetOne)
	privateRouter.GET("/getList", s.GetList)
	privateRouter.GET("/getOption", s.GetOption)
	privateRouter.POST("/updateOne", s.UpdateOne)
}

// AddOne
// @Summary		    新增单条记录
// @Description	    新增单条记录 Add Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	body		{{.PkgName}}Req.Add{{.ModelName}}Req	true	"请求参数"
// @Success		    200		{object}	response.Body{data={{.PkgName}}Model.{{.ModelName}},msg=string}
// @Router			/{{.PkgName}}/user/addOne [post]
func (s *{{.ModelName}}Service) AddOne(ctx *gin.Context) {
	var inParam {{.PkgName}}Req.Add{{.ModelName}}Req
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	createdOne, err := s.uc.Create(ctx, &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, createdOne)
}

// DelOne
// @Summary		    删除单条记录
// @Description	    删除单条记录 Del Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	query		utils.PkByInt32Param	true	"请求参数"
// @Success		    200		{object}	response.Body{msg=string}
// @Router			/{{.PkgName}}/user/delOne [get]
func (s *{{.ModelName}}Service) DelOne(ctx *gin.Context) {
	var inParam utils.PkByInt32Param
	if err := ctx.ShouldBindQuery(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Delete(ctx, inParam.Id); err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(ctx)
}

// Delmany
// @Summary		    删除多条记录
// @Description	    删除多条记录 Del Many Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	query		utils.PkByInt32sParam	true	"请求参数"
// @Success		    200		{object}	response.Body{msg=string}
// @Router			/{{.PkgName}}/user/delMany [post]
func (s *{{.ModelName}}Service) DelMany(ctx *gin.Context) {
	var inParam utils.PkByInt32sParam
	if err := ctx.ShouldBindQuery(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Delete(ctx, inParam.Id); err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(ctx)
}

// GetOne
// @Summary		    查询单条记录
// @Description	    查询单条记录 By PK Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	query		utils.PkByInt32Param	true	"请求参数"
// @Success		    200		{object}	response.Body{data={{.PkgName}}Model.{{.ModelName}},msg=string}
// @Router			/{{.PkgName}}/user/getOne [get]
func (s *{{.ModelName}}Service) GetOne(ctx *gin.Context) {
	var inParam utils.PkByInt32Param
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	createdOne, err := s.uc.Get(ctx, inParam.Id)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, createdOne)
}

// GetList
// @Summary		    查询N条记录
// @Description	    查询N条记录 List/Page Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	query		{{.PkgName}}Req.List{{.ModelName}}Req	true	"请求参数"
// @Success		    200		{object}	response.Body{data=utils.PageResult{list=[]{{.PkgName}}Model.{{.ModelName}}},msg=string}
// @Router			/{{.PkgName}}/user/getList [get]
func (s *{{.ModelName}}Service) GetList(ctx *gin.Context) {
	var inParam {{.PkgName}}Req.List{{.ModelName}}Req
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	list, total, err := s.uc.List(ctx, &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, utils.PageResult{
		List:     list,
		NotPage:  inParam.NotPage,
		Total:    total,
		PageNo:   inParam.PageNo,
		PageSize: inParam.PageSize,
	})
}

// GetOption
// @Summary		    查询N条记录
// @Description	    查询N条记录 List/Page Model To Option
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	query		{{.PkgName}}Req.List{{.ModelName}}Req	true	"请求参数"
// @Success		    200		{object}	response.Body{data=utils.PageResult{list=[]{{.PkgName}}Model.{{.ModelName}}},msg=string}
// @Router			/{{.PkgName}}/user/getOption [get]
func (s *{{.ModelName}}Service) GetOption(ctx *gin.Context) {
	var inParam {{.PkgName}}Req.List{{.ModelName}}Req
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	list, total, err := s.uc.ListByOption(ctx, &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, utils.PageResult{
		List:     list,
		NotPage:  inParam.NotPage,
		Total:    total,
		PageNo:   inParam.PageNo,
		PageSize: inParam.PageSize,
	})
}

// UpdateOne
// @Summary		    修改单条记录
// @Description	    修改单条记录 Update Model
// @Tags			{{.PkgName}}/{{.ModelName}}
// @Accept			json
// @Produce		    json
// @Param			inParam	body		{{.PkgName}}Req.Add{{.ModelName}}Req	true	"请求参数"
// @Success		    200		{object}	response.Body{data={{.PkgName}}Model.{{.ModelName}},msg=string}
// @Router			/{{.PkgName}}/user/updateOne [post]
func (s *{{.ModelName}}Service) UpdateOne(ctx *gin.Context) {
	var inParam {{.PkgName}}Req.Update{{.ModelName}}Req
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
    var fields map[string]interface{}
    mapstructure.Decode(inParam, &fields)
	updatedOne, err := s.uc.Update(ctx, inParam.Id,fields)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, updatedOne)
}
