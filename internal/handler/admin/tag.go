package admin

import (
	"log"
	"net/http"

	"github.com/taadis/blog-web/internal/pkg/mysql"
	"github.com/taadis/blog-web/internal/pkg/view"
)

type TagHandler struct {
}

func NewTagHandler() *TagHandler {
	h := new(TagHandler)
	return h
}

func (h *TagHandler) TagList(w http.ResponseWriter, r *http.Request) {
	tags, err := mysql.GetTags()
	if err != nil {
		log.Println("get ags err:", err)
		return
	}

	data := make(map[string]interface{})
	data["tags"] = tags
	view.AdminRender(data, w, "/admin/tag/list")
}
