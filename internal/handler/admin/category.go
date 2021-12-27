package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/taadis/blog-web/internal/bll"
	"github.com/taadis/blog-web/internal/pkg/view"
)

type CategoryHandler struct {
	categoryBll bll.CategoryBlli
}

func NewCategoryHandler() *CategoryHandler {
	h := new(CategoryHandler)
	h.categoryBll = bll.NewCategoryBll()
	return h
}

func (h *CategoryHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryBll.GetCategories(r.Context())
	if err != nil {
		log.Printf("CategoryHandler CategoryList categoryBll.GetCategories error:%+v", err)
		return
	}

	data := make(map[string]interface{})
	data["categories"] = categories
	view.AdminRender(data, w, "/admin/category/list")
}

func (h *CategoryHandler) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	log.Printf("category add ...")
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	category := &bll.GetCategoryResult{}
	if id > 0 {
		category, _ = h.categoryBll.GetCategory(r.Context(), &bll.GetCategoryParams{Id: id})
	}
	categories, _ := h.categoryBll.GetCategories(r.Context())

	data := make(map[string]interface{})
	data["categories"] = categories
	if category.Id > 0 {
		data["id"] = category.Id
		data["name"] = category.Name
	}
	view.AdminRender(data, w, "/admin/category/add")
}

func (h *CategoryHandler) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	err := h.categoryBll.Delete(r.Context(), id)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败,请重试"
		view.AdminRender(data, w, "401")
		return
	}

	http.Redirect(w, r, "/admin/category", http.StatusFound)
}

func (h *CategoryHandler) CategorySave(w http.ResponseWriter, r *http.Request) {
	params := new(bll.SaveCategoryParams)
	params.Id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)
	params.Name = r.FormValue("name")
	err := h.categoryBll.SaveCategory(r.Context(), params)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "保存失败,请重试"
		view.AdminRender(data, w, "401")
		return
	}

	http.Redirect(w, r, "/admin/category", http.StatusFound)
}
