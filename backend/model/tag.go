package model

import "time"

type Tags struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null; unique"`
	CreatedAt time.Time `json:"created_at"`
}

type TagsResponse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null; unique"`
	CreatedAt time.Time `json:"created_at"`
}

type PostTags struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	PostID    Posts     `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	TagID     Tags      `json:"tag_id" gorm:"foreignKey:ID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
}

type PostTagsResponse struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	PostID    Posts     `json:"post_id" gorm:"foreignKey:ID;references:ID"`
	TagID     Tags      `json:"tag_id" gorm:"foreignKey:ID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
}