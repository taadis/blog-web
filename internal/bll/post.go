package bll

import (
	"context"
	"log"
	"time"

	"github.com/taadis/blog-web/internal/pkg/model"
	"github.com/taadis/blog-web/internal/pkg/mysql"
)

type PostBlli interface {
	SavePost(ctx context.Context, params *Post) (int64, error)
	GetPost(ctx context.Context, id int64) (*Post, error)
	GetPosts(ctx context.Context, params *GetPostParams) ([]*Post, error)
	GetPostIdsByContent(ctx context.Context, content string) ([]string, error)
	DeletePost(ctx context.Context, id int64) error
	IncrView(ctx context.Context, id int64) error
}

type PostBll struct {
}

func NewPostBll() PostBlli {
	postBll := new(PostBll)
	return postBll
}

type Post struct {
	Id           int64
	Title        string
	Views        int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CategoryId   int64
	CategoryName string
	Description  string
	Content      string
	TagNames     []string
	Status       int64
}

func (b *PostBll) SavePost(ctx context.Context, params *Post) (int64, error) {
	post := new(model.Post)
	post.Id = params.Id
	post.Title = params.Title
	post.Description = params.Description
	post.Content = params.Content
	post.CategoryId = params.CategoryId
	post.Status = 1
	id, err := mysql.PostSave(post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *PostBll) GetPost(ctx context.Context, id int64) (*Post, error) {
	post, err := mysql.GetPost(id)
	if err != nil {
		log.Printf("PostBll GetPost mysql.GetPost error:%+v, id:%d", err, id)
		return nil, err
	}
	if post == nil {
		return nil, nil
	}

	result := &Post{
		Id:           post.Id,
		Title:        post.Title,
		Views:        post.Views,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		CategoryId:   post.CategoryId,
		CategoryName: post.CategoryName,
		Description:  post.Description,
		Content:      post.Content,
		TagNames:     post.TagNames,
		Status:       post.Status,
	}
	return result, nil
}

type GetPostParams struct {
	Ids        map[string][]string
	CategoryId string
	TagId      string
	PerPage    int
	Page       int
	Keyword    string
}

func (b *PostBll) GetPosts(ctx context.Context, params *GetPostParams) ([]*Post, error) {
	ps := new(mysql.PostParams)
	ps.Ids = params.Ids
	ps.CategoryId = params.CategoryId
	ps.TagId = params.TagId
	ps.PerPage = params.PerPage
	ps.Page = params.Page
	ps.Keyword = params.Keyword
	posts, err := mysql.GetPosts(ps)
	if err != nil {
		log.Printf("PostBll GetPosts mtsql.GetPosts error:%+v", err)
		return nil, err
	}
	if posts == nil {
		return nil, nil
	}

	result := make([]*Post, 0, len(posts))
	for _, v := range posts {
		result = append(result, &Post{
			Id:           v.Id,
			Title:        v.Title,
			Views:        v.Views,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			CategoryId:   v.CategoryId,
			CategoryName: v.CategoryName,
			Description:  v.Description,
			Content:      v.Content,
			TagNames:     v.TagNames,
			Status:       v.Status,
		})
	}
	return result, nil
}

func (b *PostBll) GetPostIdsByContent(ctx context.Context, content string) ([]string, error) {
	postIds, err := mysql.GetPostIdsByContent(content)
	if err != nil {
		log.Printf("PostBll GetPostIdsByContent mysql.GetPostIdsByContent error:%+v, content:%s", err, content)
		return nil, err
	}

	return postIds, nil
}

func (b *PostBll) DeletePost(ctx context.Context, id int64) error {
	err := mysql.PostDelete(id)
	if err != nil {
		log.Printf("PostBll DeletePost mysql.PostDelete(%d) error:%+v", id, err)
		return err
	}

	//if !conf.Conf.Elasticsearch.Disable {
	//	go es.DeletePost(es.Post{Id: id})
	//}

	return nil
}

func (b *PostBll) IncrView(ctx context.Context, id int64) error {
	err := mysql.IncrView(id)
	if err != nil {
		log.Printf("PostBll IncrView mysql.IncrView error:%+v, id:%d", err, id)
		return err
	}

	return nil
}
