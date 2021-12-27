package mysql

import (
	"log"

	"github.com/taadis/blog-web/internal/pkg/model"
)

func GetUser(email string) (*model.User, error) {
	user := new(model.User)
	row := db.QueryRow("select id, email, password from user where email=?", email)
	err := row.Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		log.Printf("mysql GetUser row.Scan error:%+v", err)
		return nil, err
	}

	return user, nil
}

func AddUser(user *model.User) (int64, error) {
	rs, err := db.Exec("insert into user (email, password) values (?, ?)", user.Email, user.Password)
	if err != nil {
		log.Printf("AddUser db.Exec error:%+v, user data:%+v", err, user)
		return 0, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}
