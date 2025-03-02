package model

import (
	"time"
)

type Experiences struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string    `json:"userId" gorm:"uniqueIndex;not null;references:id"`
	Work        string    `json:"work" gorm:"not null;"`
	Skills      string    `json:"skills" gorm:"not null;"`
	SelfPR      string    `json:"selfPR" gorm:"not null;"`
	FutureGoals string    `json:"futureGoals" gorm:"not null;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"not null"`
	User        Users     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type InputExperience struct {
	Work        string
	Skills      string
	SelfPR      string
	FutureGoals string
}
