package model

// User reprensents a registered user
type User struct {
	ID    uint `gorm:"primary_key"`
	Token string
	Admin bool
}
