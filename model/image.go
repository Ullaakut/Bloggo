package model

import (
	"time"
)

// Image reprensents an uploaded image file
type Image struct {
	ID        uint      `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Data      []byte    `json:"data" validate:"required"`
	Size      int64     `json:"size,omitempty"`
}
