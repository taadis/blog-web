package model

type Tag struct {
	Id        int64
	Name      string
	Count     int64
	CreatedAt string
	UpdatedAt string
}

type TagPost struct {
	Id     int64
	TagId  int64
	PostId int64
}
