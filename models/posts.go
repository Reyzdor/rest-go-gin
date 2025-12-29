package models

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	UserID    int
	Username  string
	CreatedAt time.Time
	Price     int
	MainImage string
	Images    []string
}

type PostImage struct {
	ID        int
	PostID    int
	ImagePath string
	IsMain    bool
	SortOrder int
}
