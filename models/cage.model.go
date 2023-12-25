package models

import (
	"time"
)

type Cage struct {
	Id        int       `json:"id" gorm:"unique"`
	Status    bool      `gorm:"not null" json:"status"`
	Capacity  int       `gorm:"not null" json:"capacity"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

type CageRequest struct {
	Status    bool      `json:"status"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
