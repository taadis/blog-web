package mysql

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/taadis/blog-web/internal/pkg/model"
)

type PostParams struct {
	Ids        map[string][]string
	CategoryId string
	TagId      string
	PerPage    int
	Page       int
	Keyword    string
}

func GetPosts(params *PostParams) (posts []*model.Post, err error) {
	var condition []string
	var args []interface{}

	if len(params.CategoryId) > 0 {
		condition = append(condition, "category_id=?")
		args = append(args, params.CategoryId)
	}

	if len(params.TagId) > 0 {
		condition = append(condition, "JSON_CONTAINS(tag_ids,?)")
		args = append(args, params.TagId)
	}
	if len(params.Keyword) > 0 {
		keywordSql := "title like ? or description like ?"
		args = append(args, "%"+params.Keyword+"%", "%"+params.Keyword+"%")
		if len(params.Ids["ids"]) > 0 {
			keywordSql += " or id in (?)"
			args = append(args, strings.Join(params.Ids["ids"], ","))
		}
		if len(params.Ids["category_ids"]) > 0 {
			keywordSql += " or category_id in (?)"
			args = append(args, strings.Join(params.Ids["category_ids"], ","))
		}
		if len(params.Ids["tag_ids"]) > 0 {
			for _, tagId := range params.Ids["tag_ids"] {
				keywordSql += " or JSON_CONTAINS(tag_ids,?)"
				args = append(args, tagId)
			}
		}
		condition = append(condition, "("+keywordSql+")")
	}

	querySql := "select id,title,description,created_at,updated_at,category_id,views from post"

	condition = append(condition, "status=1")
	if len(condition) > 0 {
		querySql += ` where ` + strings.Join(condition, " and ")
	}

	querySql += " order by is_top desc,updated_at desc"

	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		rows.Scan(&post.Id, &post.Title, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &post.Views, &post.Description)
		//json.Unmarshal(tags, &post.TagIds)
		posts = append(posts, &post)
	}
	return
}

func GetPost(id int64) (*model.Post, error) {
	//var tags []byte
	var post model.Post
	row := db.QueryRow("select id,title,description,content,created_at,updated_at,category_id,views from post where id=? limit 1", id)
	err := row.Scan(&post.Id, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &post.Views)
	if err != nil {
		log.Printf("GetPost row.Scan error:%+v", err)
		return nil, err
	}
	//json.Unmarshal(tags, &post.TagIds)
	return &post, nil
}

func IncrView(id int64) error {
	_, err := db.Exec("update post set views=views+1 where id = ?", id)
	if err != nil {
		log.Printf("update post views err id:%s, err%v", id, err)
		return err
	}

	return nil
}

func PostDelete(id int64) error {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalln("post db conn err ", err)
		return err
	}

	_, err = conn.Exec("delete from post where id=?", id)
	if err != nil {
		log.Printf("post %d delete err %v", id, err)
		conn.Rollback()
		return err
	}

	err = conn.Commit()
	if err != nil {
		log.Printf("post delete conn.Commit error:%+v, id:%d", err, id)
		return err
	}

	return nil
}

func PostSave(post *model.Post) (int64, error) {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalf("PostSave db.Begin error:%+v", err)
		return 0, err
	}

	//var tagIds []byte
	//tagIds, _ = json.Marshal(post.TagIds)
	var rs sql.Result
	if post.Id > 0 {
		_, err = conn.Exec("update post set title=?,description=?,content=?,category_id=? where id=?",
			post.Title,
			post.Description,
			post.Content,
			post.CategoryId,
			post.Id,
		)
		if err != nil {
			log.Printf("PostSave update conn.Exec error:%+v", err)
			conn.Rollback()
			return 0, err
		}
	} else {
		rs, err = conn.Exec("insert into post (title,description,content,category_id,status) values (?,?,?,?,?)",
			post.Title,
			post.Description,
			post.Content,
			post.CategoryId,
			post.Status,
		)
		if err != nil {
			log.Printf("PostSave conn.Exec insert error:%+v", err)
			conn.Rollback()
			return 0, err
		}

		post.Id, err = rs.LastInsertId()
	}
	conn.Commit()
	return post.Id, nil
}

func GetPostCountByTagId(id string) (int, error) {
	var count *int
	row := db.QueryRow("select count(*) from post where JSON_CONTAINS(tag_ids,'" + id + "')")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return *count, err
}

func GetPostIdsByContent(content string) ([]string, error) {
	rows, err := db.Query("select id from post where content like ?", "%"+content+"%")
	if err != nil {
		return nil, err
	}
	var postIds []string
	defer rows.Close()
	for rows.Next() {
		var postId string
		rows.Scan(&postId)
		postIds = append(postIds, postId)
	}
	return postIds, nil
}
