package bll

import (
	"context"
	"log"

	"github.com/taadis/blog-web/internal/pkg/model"
	"github.com/taadis/blog-web/internal/pkg/mysql"
)

type CategoryBlli interface {
	SaveCategory(ctx context.Context, params *SaveCategoryParams) error
	GetCategory(ctx context.Context, params *GetCategoryParams) (*GetCategoryResult, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategoryIdsByName(ctx context.Context, name string) ([]string, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryBll struct {
}

func NewCategoryBll() CategoryBlli {
	bll := new(CategoryBll)
	return bll
}

func (b *CategoryBll) GetCategoryIdsByName(ctx context.Context, name string) ([]string, error) {
	categoryIds, err := mysql.GetCategoryIdsByName(name)
	if err != nil {
		return nil, err
	}

	return categoryIds, err
}

func (b *CategoryBll) GetCategories(ctx context.Context) ([]*Category, error) {
	categories, err := mysql.GetCategories()
	if err != nil {
		log.Printf("CategoryBll GetCategories mysql.GetCategories error:%+v", err)
		return nil, err
	}

	result := make([]*Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, &Category{
			Id:   category.Id,
			Name: category.Name,
		})
	}
	return result, nil
}

func (b *CategoryBll) GetCategory(ctx context.Context, params *GetCategoryParams) (*GetCategoryResult, error) {
	category := mysql.GetCategory(params.Id)
	if category == nil {
		return nil, nil
	}

	result := &GetCategoryResult{
		Id:   category.Id,
		Name: category.Name,
	}
	return result, nil
}

type GetCategoryParams struct {
	Id int64
}

type GetCategoryResult struct {
	Id   int64
	Name string
}

func (b *CategoryBll) SaveCategory(ctx context.Context, params *SaveCategoryParams) error {
	category := &model.Category{
		Id:   params.Id,
		Name: params.Name,
	}
	_, err := mysql.CategorySave(category)
	if err != nil {
		log.Printf("CategoryBll Save mysql.CategorySave error:%+v,params:%+v", err, params)
		return err
	}

	return nil
}

type SaveCategoryParams struct {
	Category
}

type Category struct {
	Id        int64
	Name      string
	CreatedAt string
	UpdatedAt string
	Cur       int64
}

func (b *CategoryBll) Delete(ctx context.Context, id int64) error {
	err := mysql.CategoryDelete(id)
	if err != nil {
		log.Printf("CategoryBll Delete mysql.CategoryDelete error:%+v, id:%d", err, id)
		return err
	}

	return nil
}
