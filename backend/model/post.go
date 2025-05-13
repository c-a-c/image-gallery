package model

import "time"

type Posts struct {
	ID 	       uint      `json:"id" gorm:"primary_key"`
	User       Users     `json:"user" gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	ImageUrl   string    `json:"image_url"`
	Caption    string    `json:"caption"`
	CreatedAt  time.Time `json:"datetime"`
	UpdateAt   time.Time `json:"update_at"`
}

type PostsResponse struct {
	ID 	       uint      `json:"id" gorm:"primary_key"`
	ImageUrl   string    `json:"image_url"`
	Caption    string    `json:"caption"`
	CreatedAt  time.Time `json:"datetime"`
	UpdateAt   time.Time `json:"update_at"`
}