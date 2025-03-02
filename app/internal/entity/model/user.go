package model

import (
	"time"
)

type Users struct {
	ID        string    `json:"id" gorm:"primaryKey;unique;not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
}
