package model

import (
	"time"
)

// BlogPost reprensents a blog post
type BlogPost struct {
	ID        uint      `json:"id,omitempty" gorm:"primary_key"`
	Author    string    `json:"author,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt time.Time // TODO maybe if I have time
	// Comments   []Comment // TODO maybe if I have time
	// Views   uint // TODO maybe if I have time
}

// type Comment struct {
// 	Author      string
// 	Content     string
// 	CreatedDate time.Time
// }
