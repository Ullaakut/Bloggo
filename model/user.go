package model

// User reprensents a registered user
type User struct {
	ID          uint `gorm:"primary_key"`
	TokenUserID string
	Email       string `validate:"required,email"`
	Password    string `validate:"required,min=10"`
	IsAdmin     bool   `json:"is_admin"`
}
