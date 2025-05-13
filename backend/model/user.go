// user.go

package model

import "time"

type Users struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email	  string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	AvatarUrl string	`json:"avatar_url" gorm:"not null"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdateAt  time.Time `json:"update_at" gorm:"not null"`
	IsMember  bool      `json:"is_member" gorm:"not null"`
}
