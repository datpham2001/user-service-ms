package entity

import "time"

type User struct {
	ID           string    `gorm:"primary_key"`
	Email        string    `gorm:"unique"`
	Password     string    `gorm:"not null"`
	Username     string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
