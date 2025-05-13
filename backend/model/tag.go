// tag.go

package model

import "time"

type Tags struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null; unique"`
	CreatedAt time.Time `json:"created_at"`
}
