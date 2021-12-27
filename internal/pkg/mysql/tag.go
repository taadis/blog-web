package mysql

import (
	"log"

	"github.com/taadis/blog-web/internal/pkg/model"
	logger "github.com/taadis/blog-web/pkg/log"
	"go.uber.org/zap"
)

func GetTags() (tags []*model.Tag, err error) {
	rows, err := db.Query("select id,name,count from tag order by count desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag model.Tag
		rows.Scan(&tag.Id, &tag.Name, &tag.Count)
		tags = append(tags, &tag)
	}
	return
}

func GetTagIdsByName(name string) (tagIds []string, err error) {
	rs, err := db.Query("select id from tag where name like ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var tagId string
		rs.Scan(&tagId)
		tagIds = append(tagIds, tagId)
	}
	return
}

func AddTag(tag *model.Tag) (int64, error) {
	rs, err := db.Exec("insert into tag (name) values (?)", tag.Name)
	if err != nil {
		log.Printf("AddTag db.Exec error:%+v", err)
		return 0, err
	}
	id, err := rs.LastInsertId()
	return id, err
}

func IncrTagCount(id string) {
	count, err := GetPostCountByTagId(id)
	if err != nil {
		logger.Error("incr_tag_count_err", zap.Error(err))
		return
	}

	db.Exec("update tag set count=? where id = ?", count, id)
}
