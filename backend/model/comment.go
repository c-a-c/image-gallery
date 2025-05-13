package model

import "time"

type Comments struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	PostID       Posts     `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	UserID       Users     `json:"user_id" gorm:"foreignKey:ID;references:ID"`
	Context      string    `json:"context"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt	 time.Time `json:"update_at"`
}

type CommentsResponse struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	PostID       Posts     `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	UserID       Users     `json:"user_id" gorm:"foreignKey:ID;references:ID"`
	Context      string    `json:"context"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt	 time.Time `json:"update_at"`
}
