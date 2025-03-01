package entity

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"gorm:unique not null"`
	Name      string    `json:"name" gorm:"unique not null"`
	Email     string    `json:"email" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
