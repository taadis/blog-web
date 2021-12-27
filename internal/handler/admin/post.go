package admin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/taadis/blog-web/conf"
	"github.com/taadis/blog-web/internal/bll"
	"github.com/taadis/blog-web/internal/pkg/utils"
	"github.com/taadis/blog-web/internal/pkg/view"
	//"github.com/taadis/blog-web/internal/pkg/es"
)

type PostHandler struct {
	postBll     bll.PostBlli
	categoryBll bll.CategoryBlli
	tagBll      bll.TagBlli
}

func NewPostHandler() *PostHandler {
	h := new(PostHandler)
	h.postBll = bll.NewPostBll()
	h.categoryBll = bll.NewCategoryBll()
	h.tagBll = bll.NewTagBll()
	return h
}

func (h *PostHandler) PostList(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")
	tagId := r.URL.Query().Get("tag_id")
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
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
	posts, err := h.postBll.GetPosts(r.Context(), &bll.GetPostParams{
		CategoryId: categoryId,
		TagId:      tagId,
		PerPage:    perPage,
		Page:       page,
	})
	if err != nil {
		log.Printf("PostHandler PostList postBll.GetPost error:%+v", err)
		return
	}

	categories, err := h.categoryBll.GetCategories(r.Context())
	if err != nil {
		fmt.Println("get categories err:", err)
		return
	}

	categoryMap := make(map[int64]*bll.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}
	for index, post := range posts {
		if category, ok := categoryMap[post.CategoryId]; ok {
			posts[index].CategoryName = category.Name
		}
		//posts[index].CategoryName = categoryMap[post.CategoryId].Name
	}
	data := make(map[string]interface{})
	data["posts"] = posts
	data["categories"] = categories
	data["page"] = page
	data["pre_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(prePage))
	data["next_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(nextPage))
	view.AdminRender(data, w, "/admin/post/list")
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

func (h *PostHandler) PostAdd(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	post := new(bll.Post)
	if id > 0 {
		post, _ = h.postBll.GetPost(r.Context(), id)
	}
	categories, _ := h.categoryBll.GetCategories(r.Context())
	data["categories"] = categories

	if post.Id > 0 {
		for i := range categories {
			categories[i].Cur = post.CategoryId
		}
		data["id"] = post.Id
		data["title"] = post.Title
		data["description"] = post.Description
		data["content"] = post.Content
		data["category_id"] = post.CategoryId
		//data["tag_ids"] = post.TagIds
		//data["tags"] = getTags(post)
	}
	view.AdminRender(data, w, "/admin/post/add")
}

// func (h *PostHandler) getTags(post *model.Post) string {
// 	var tags []string
// 	if len(post.TagIds) > 0 {
// 		allTags, _ := h.tagBll.GetTags()

// 		tagsById := make(map[int64]model.Tag)
// 		for _, tag := range allTags {
// 			tagsById[tag.Id] = tag
// 		}
// 		for _, tagId := range post.TagIds {
// 			tags = append(tags, tagsById[tagId].Name)
// 		}

// 	}
// 	return strings.Join(tags, ",")
// }

func (h *PostHandler) PostDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	err := h.postBll.DeletePost(r.Context(), id)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败,请重试"
		view.AdminRender(data, w, "401")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *PostHandler) PostSave(w http.ResponseWriter, r *http.Request) {
	post := new(bll.Post)
	post.Id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Content = r.FormValue("content")
	post.CategoryId, _ = strconv.ParseInt(r.FormValue("category"), 10, 64)
	//tags := r.FormValue("tags")
	post.Status = 1
	_, err := h.postBll.SavePost(r.Context(), post)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		view.AdminRender(data, w, "401")
		return
	}
	//if !conf.Conf.Elasticsearch.Disable {
	//	category := h.categoryBll.GetCategory(post.CategoryId)
	//	go es.SavePost(es.Post{
	//		Id:          post.Id,
	//		Title:       post.Title,
	//		Description: post.Description,
	//		Content:     post.Content,
	//		Tags:        tags,
	//		Category:    category.Name,
	//	})
	//}
	//for _, tagId := range post.TagIds {
	//	go mysql.IncrTagCount(strconv.Itoa(tagId))
	//}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *PostHandler) getTagIds(ctx context.Context, tags string) (tagIds []int64) {
	tagNames := strings.Split(tags, ",")
	tagNames = utils.RemoveDuplicateElement(tagNames)
	allTags, _ := h.tagBll.GetTags(ctx)
	var allTagNames []string
	allTagByName := make(map[string]*bll.Tag)
	for _, tag := range allTags {
		allTagNames = append(allTagNames, tag.Name)
		allTagByName[tag.Name] = tag
	}
	for _, tagName := range tagNames {
		if utils.StrInArray(tagName, allTagNames) {
			tagIds = append(tagIds, allTagByName[tagName].Id)
		} else {
			var newTag bll.Tag
			newTag.Name = tagName
			newTagId, _ := h.tagBll.AddTag(context.TODO(), &bll.Tag{Name: tagName})
			if newTagId > 0 {
				tagIds = append(tagIds, newTagId)
			}
		}
	}
	return
}
