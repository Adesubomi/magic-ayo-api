package auth

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
