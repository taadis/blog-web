package admin

import (
	"log"
	"net/http"

	"github.com/taadis/blog-web/internal/bll"
	"github.com/taadis/blog-web/internal/pkg/view"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userBll bll.UserBlli
}

func NewAuthHandler() *AuthHandler {
	h := new(AuthHandler)
	h.userBll = bll.NewUserBll()
	return h
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	view.AdminRender(data, w, "/admin/login")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	view.AdminRender(data, w, "/admin/register")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("logout...")
	cookie := &http.Cookie{
		Name:   "email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	repassword := r.FormValue("repassword")
	if email == "" || password == "" || repassword == "" || password != repassword {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	addUserParams := new(bll.AddUserParams)
	addUserParams.Email = email
	addUserParams.Password = string(hashPassword)
	// 不允许注册
	return
	_, err := h.userBll.AddUser(r.Context(), addUserParams)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:  "email",
		Value: email,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, _ := h.userBll.GetUser(r.Context(), email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := make(map[string]interface{})
		data["msg"] = "密码不正确,请重试"
		view.AdminRender(data, w, "/admin/401")
		return
	}
	cookie := &http.Cookie{
		Name:  "email",
		Value: email,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin", http.StatusFound)
}
