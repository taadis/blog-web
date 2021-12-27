package bll

import (
	"context"
	"log"

	"github.com/taadis/blog-web/internal/pkg/model"
	"github.com/taadis/blog-web/internal/pkg/mysql"
)

type UserBlli interface {
	AddUser(ctx context.Context, params *AddUserParams) (int64, error)
	GetUser(ctx context.Context, email string) (*User, error)
}

type UserBll struct {
}

func NewUserBll() UserBlli {
	bll := new(UserBll)
	return bll
}

type User struct {
	Id       int64
	Email    string
	Password string
}

type AddUserParams struct {
	User
}

func (b *UserBll) AddUser(ctx context.Context, params *AddUserParams) (int64, error) {
	userModel := &model.User{
		Id:       params.Id,
		Email:    params.Email,
		Password: params.Password,
	}
	id, err := mysql.AddUser(userModel)
	if err != nil {
		log.Printf("UserBll AddUser mysql.AddUser error:%+v", err)
		return 0, err
	}

	return id, nil
}

func (b *UserBll) GetUser(ctx context.Context, email string) (*User, error) {
	u, err := mysql.GetUser(email)
	if err != nil {
		log.Printf("UserBll GeUser mysql.GetUser error:%+v", err)
		return nil, err
	}
	if u == nil {
		return nil, nil
	}

	user := &User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
	return user, nil
}
