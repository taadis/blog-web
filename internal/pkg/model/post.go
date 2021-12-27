package model

import "time"

type Post struct {
	Id           int64
	Title        string
	Views        int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CategoryId   int64
	CategoryName string
	Description  string
	Content      string
	TagNames     []string
	Status       int64
}
