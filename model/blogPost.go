package model

import (
	"time"
)

// BlogPost reprensents a blog post
type BlogPost struct {
	ID        uint      `json:"id,omitempty" gorm:"primary_key"`
	Author    string    `json:"author,omitempty"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Comments   []Comment // TODO maybe if I have time
	// Views   uint // TODO maybe if I have time
}

// type Comment struct {
// 	Author      string
// 	Content     string
// 	CreatedDate time.Time
// }
