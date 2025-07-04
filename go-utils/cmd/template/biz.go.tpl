package {{.PkgName}}

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go-utils/utils"
	{{.PkgName}}Model "go-kgin-platform/internal/dal/model/{{.PkgName}}"
	{{.PkgName}}Repo "go-kgin-platform/internal/dal/repo/{{.PkgName}}"
	{{.PkgName}}Req "go-kgin-platform/internal/dal/req/{{.PkgName}}"
	{{.PkgName}}Vo "go-kgin-platform/internal/dal/vo/{{.PkgName}}"
)

var (
	Err{{.ModelName}}Existed = errors.New("名称重复")
)

// {{.ModelName}}UseCase is a {{.ModelName}} usecase.
type {{.ModelName}}UseCase struct {
	repo {{.PkgName}}Repo.{{.ModelName}}Repo
	log  *log.Helper
}

// New{{.ModelName}}UseCase new a {{.ModelName}} usecase.
func New{{.ModelName}}UseCase(
	repo {{.PkgName}}Repo.{{.ModelName}}Repo,
	logger log.Logger,
) *{{.ModelName}}UseCase {
	return &{{.ModelName}}UseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/{{.ModelName}}")),
	}
}

func (uc *{{.ModelName}}UseCase) beforeCheck(ctx context.Context, req *{{.PkgName}}Req.Add{{.ModelName}}Req) error {
	old, err := uc.repo.SelectOne(ctx, &{{.PkgName}}Model.{{.ModelName}}{

	})
	if err != nil {
		return err
	}
	if old != nil {
		// return errors.Wrapf(Err{{.ModelName}}Existed, "请检查[%s]", req.UserName)
	}
	return nil
}

func (uc *{{.ModelName}}UseCase) Create(ctx context.Context, req *{{.PkgName}}Req.Add{{.ModelName}}Req) (*{{.PkgName}}Model.{{.ModelName}}, error) {
	if err := uc.beforeCheck(ctx, req); err != nil {
		return nil, err
	}
	{{.SnakeModelName}}Model := {{.PkgName}}Model.{{.ModelName}}{}
	err := uc.repo.Insert(ctx, &{{.SnakeModelName}}Model)
	return &{{.SnakeModelName}}Model, err
}

func (uc *{{.ModelName}}UseCase) Delete(ctx context.Context, id interface{}) error {
	_, err := uc.repo.DeleteByPK(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *{{.ModelName}}UseCase) Get(ctx context.Context, id interface{}) (*{{.PkgName}}Vo.{{.ModelName}}Vo, error) {
	one, err := uc.repo.SelectOneByPK(ctx, id)
	if err != nil {
		return nil, err
	}
	return &{{.PkgName}}Vo.{{.ModelName}}Vo{
		{{.ModelName}}: one,
	}, nil
}

func (uc *{{.ModelName}}UseCase) ListByOption(ctx context.Context, req *{{.PkgName}}Req.List{{.ModelName}}Req) ([]*utils.Option, int64, error) {
	query := uc.repo.DB(ctx).Model(&{{.PkgName}}Model.{{.ModelName}}{})
	query = query.Select("id as value, name as label")
	var res []*utils.Option
	total, err := uc.repo.ListPageByOption(ctx, query, req.PageParam, &res)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		total = int64(len(res))
	}
	return res, total, nil
}

func (uc *{{.ModelName}}UseCase) List(ctx context.Context, req *{{.PkgName}}Req.List{{.ModelName}}Req) ([]*{{.PkgName}}Model.{{.ModelName}}, int64, error) {
	query := uc.repo.DB(ctx).Model(&{{.PkgName}}Model.{{.ModelName}}{})

	return uc.repo.ListPage(ctx, query, req.PageParam)
}

func (uc *{{.ModelName}}UseCase) Update(ctx context.Context, id interface{}, field map[string]interface{}) (*{{.PkgName}}Vo.{{.ModelName}}Vo, error) {
	row, err := uc.repo.UpdateByPKWithMap(ctx, id, field)
	if err != nil {
		return nil, err
	}
	if row == 0 {
		return nil, Err{{.ModelName}}Existed
	}
	return uc.Get(ctx, id)
}
