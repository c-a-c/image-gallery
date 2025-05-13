package model

import "time"

type Likes struct {
	ID          uint `json:"id" gorm:"primary_key"`
	PostID      Posts `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	UserID      Users `json:"user_id" gorm:"foreignKey:ID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
}

type LikesResponse struct {
	ID          uint `json:"id" gorm:"primary_key"`
	PostID      Posts `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	UserID      Users `json:"user_id" gorm:"foreignKey:ID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
}