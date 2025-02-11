package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Username  string    `gorm:"type:varchar(255);not null;unique"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
