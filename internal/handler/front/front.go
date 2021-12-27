package front

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/taadis/blog-web/conf"
	"github.com/taadis/blog-web/internal/bll"
	"github.com/taadis/blog-web/internal/pkg/model"
	"github.com/taadis/blog-web/internal/pkg/view"
)

type FrontHandler struct {
	categoryBll bll.CategoryBlli
	postBll     bll.PostBlli
	tagBll      bll.TagBlli
}

func NewFrontHandler() *FrontHandler {
	h := new(FrontHandler)
	h.categoryBll = bll.NewCategoryBll()
	h.postBll = bll.NewPostBll()
	h.tagBll = bll.NewTagBll()
	return h
}

func (h *FrontHandler) Index(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")
	tagId := r.URL.Query().Get("tag_id")
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	keyword := r.URL.Query().Get("keyword")
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 1 {
		page = 1
	}
	prePage := page - 1
	nextPage := page + 1
	if prePage <= 1 {
		prePage = 1
	}
	params := &bll.GetPostParams{
		CategoryId: categoryId,
		TagId:      tagId,
		PerPage:    perPage,
		Page:       page,
	}
	// 关键词搜索：标题、描述、内容、分类、标签等
	// todo es search
	if len(keyword) > 0 {
		params.Keyword = keyword
		params.Ids = make(map[string][]string)
		postIds, err := h.postBll.GetPostIdsByContent(r.Context(), keyword)
		if err == nil {
			params.Ids["ids"] = postIds
		}
		categoryIds, err := h.categoryBll.GetCategoryIdsByName(r.Context(), keyword)
		if err == nil {
			params.Ids["category_ids"] = categoryIds
		}
		tagIds, err := h.tagBll.GetTagIdsByName(r.Context(), keyword)
		if err == nil {
			params.Ids["tag_ids"] = tagIds
		}
	}

	posts, err := h.postBll.GetPosts(r.Context(), params)
	if err != nil {
		log.Printf("Index postBll.GetPosts error:%+v", err)
		return
	}

	categories, err := h.categoryBll.GetCategories(r.Context())
	if err != nil {
		log.Printf("FrontHandler Index categoryBll.GetCategories error:%+v", err)
		return
	}

	categoryMap := make(map[int64]*bll.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}
	for _, post := range posts {
		if category, ok := categoryMap[post.CategoryId]; ok {
			post.CategoryName = category.Name
		}
	}
	//allTags, _ := h.bllTag.GetTags(r.Context())
	data := make(map[string]interface{})
	data["posts"] = posts
	data["categories"] = categories
	data["page"] = page
	data["pre_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(prePage))
	data["next_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(nextPage))
	view.Render(data, w, "index")
}

func getPageUrl(categoryId string, tagId string, page string) string {
	var params []string
	if len(categoryId) > 0 {
		params = append(params, "category_id="+categoryId)
	}
	if len(tagId) > 0 {
		params = append(params, "tag_id="+tagId)
	}
	params = append(params, "page="+page)
	return conf.Conf.App.Host + "?" + strings.Join(params, "&")
}

func (h *FrontHandler) PostInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	post, _ := h.postBll.GetPost(r.Context(), id)
	category, _ := h.categoryBll.GetCategory(r.Context(), &bll.GetCategoryParams{Id: post.CategoryId})
	//category := h.categoryBll.GetCategory(post.CategoryId)
	//allTags, _ := h.tagBll.GetTags()
	//tagIds := post.TagIds
	//tagsById := make(map[int]model.Tag)
	//for _, tag := range allTags {
	//	tagsById[tag.Id] = tag
	//}
	var tags []model.Tag
	//for _, tagId := range tagIds {
	//	tags = append(tags, tagsById[tagId])
	//}
	_ = h.postBll.IncrView(r.Context(), id)
	post.CategoryName = category.Name
	data := make(map[string]interface{})
	data["post"] = post
	data["tags"] = tags
	data["title"] = post.Title
	data["description"] = post.Description
	view.Render(data, w, "post")
}

func (h *FrontHandler) Tag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tags, _ := bll.NewTagBll().GetTags(ctx)
	data := make(map[string]interface{})
	data["title"] = "标签"
	data["description"] = "佛语的标签"
	data["tags"] = tags
	view.Render(data, w, "tag")
}

func (h *FrontHandler) Page(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	//page := h.pageBll.GetPage(id)
	page, _ := h.postBll.GetPost(r.Context(), id)
	data := make(map[string]interface{})
	data["title"] = page.Title
	data["description"] = page.Description
	data["content"] = page.Content
	data["page"] = page
	view.Render(data, w, "page")
}
