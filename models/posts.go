package models

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	UserID    int
	Username  string
	CreatedAt time.Time
	Image     string
	Price     int
}
