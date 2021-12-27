package bll

import (
	"context"
	"log"

	"github.com/taadis/blog-web/internal/pkg/model"
	"github.com/taadis/blog-web/internal/pkg/mysql"
)

type TagBlli interface {
	AddTag(ctx context.Context, params *Tag) (int64, error)
	GetTags(ctx context.Context) ([]*Tag, error)
	GetTagIdsByName(ctx context.Context, name string) (tagIds []string, err error)
}

type TagBll struct {
}

func NewTagBll() TagBlli {
	tagBll := new(TagBll)
	return tagBll
}

func (b *TagBll) GetTagIdsByName(ctx context.Context, name string) ([]string, error) {
	tagIds, err := mysql.GetTagIdsByName(name)
	if err != nil {
		log.Printf("TagBll GetTagIdsByName mysql.GetTagIdsByName error:%+v, name:%s", err, name)
		return nil, err
	}

	return tagIds, nil
}

func (b *TagBll) AddTag(ctx context.Context, params *Tag) (int64, error) {
	newTag := new(model.Tag)
	newTag.Name = params.Name
	id, err := mysql.AddTag(newTag)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *TagBll) GetTags(ctx context.Context) ([]*Tag, error) {
	tags, err := mysql.GetTags()
	if err != nil {
		log.Printf("GetTags mysql.GetTags error:%+v", err)
		return nil, err
	}

	tagList := make([]*Tag, 0)
	for _, tag := range tags {
		tagList = append(tagList, &Tag{
			Id:        tag.Id,
			Name:      tag.Name,
			Count:     tag.Count,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	return tagList, nil
}

type Tag struct {
	Id        int64
	Name      string
	Count     int64
	CreatedAt string
	UpdatedAt string
}
