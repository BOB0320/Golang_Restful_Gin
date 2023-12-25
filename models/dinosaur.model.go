package models

import (
	"time"
)

type Dinosaur struct {
	Id        int    `json:"id" gorm:"unique"`
	Name      string `gorm:"type:varchar(255);not null"`
	Type      string `gorm:"not null"`
	Spec      string `gorm:"not null"`
	CageId    int    `gorm:"not null" json:"cage_id"`
	Cage      Cage
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type DinosaurRequest struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Spec      string    `json:"spec"`
	CageId    int       `json:"cage_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

var DinosaurTypeList = map[string][]string{
	"herbivore": {"Brachiosaurus", "Stegosaurus", "Ankylosaurus", "Triceratops"},
	"carnivore": {"Tyrannosaurus", "Velociraptor", "Spinosaurus", "Megalosaurus"},
}
