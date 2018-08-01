package model

import (
	"time"
)

// BlogPost reprensents a blog post
type BlogPost struct {
	Author      string
	Title       string
	Content     string
	CreatedDate time.Time
	// UpdatedDate time.Time // TODO maybe if I have time
	// Comments   []Comment // TODO maybe if I have time
	// Views   uint // TODO maybe if I have time
}

// type Comment struct {
// 	Author      string
// 	Content     string
// 	CreatedDate time.Time
// }
