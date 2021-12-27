package routers

import (
	"net/http"

	"github.com/taadis/blog-web/internal/handler/admin"
	"github.com/taadis/blog-web/internal/handler/front"
	"github.com/taadis/blog-web/internal/routers/middleware"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func InitRouter() *http.ServeMux {
	mux := &http.ServeMux{}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	frontHandler := front.NewFrontHandler()
	mux.HandleFunc("/", frontHandler.Index)
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/post", frontHandler.PostInfo)
	mux.HandleFunc("/page", frontHandler.Page)
	mux.HandleFunc("/tag", frontHandler.Tag)

	authHandler := admin.NewAuthHandler()
	mux.HandleFunc("/admin/login", authHandler.Login)
	mux.HandleFunc("/admin/register", authHandler.Register)
	mux.HandleFunc("/admin/signin", authHandler.Signin)
	mux.HandleFunc("/admin/signup", authHandler.Signup)
	mux.Handle("/admin/logout", middleware.AuthWithCookie(http.HandlerFunc(authHandler.Logout)))

	postHandler := admin.NewPostHandler()
	mux.Handle("/admin/", middleware.AuthWithCookie(http.HandlerFunc(postHandler.PostList)))
	mux.Handle("/admin/post/add", middleware.AuthWithCookie(http.HandlerFunc(postHandler.PostAdd)))
	mux.Handle("/admin/post/save", middleware.AuthWithCookie(http.HandlerFunc(postHandler.PostSave)))
	mux.Handle("/admin/post/delete", middleware.AuthWithCookie(http.HandlerFunc(postHandler.PostDelete)))

	categoryHandler := admin.NewCategoryHandler()
	mux.Handle("/admin/category", middleware.AuthWithCookie(http.HandlerFunc(categoryHandler.CategoryList)))
	mux.Handle("/admin/category/add", middleware.AuthWithCookie(http.HandlerFunc(categoryHandler.CategoryAdd)))
	mux.Handle("/admin/category/save", middleware.AuthWithCookie(http.HandlerFunc(categoryHandler.CategorySave)))
	mux.Handle("/admin/category/delete", middleware.AuthWithCookie(http.HandlerFunc(categoryHandler.CategoryDelete)))

	tagHandler := admin.NewTagHandler()
	mux.Handle("/admin/tag", middleware.AuthWithCookie(http.HandlerFunc(tagHandler.TagList)))

	return mux
}
