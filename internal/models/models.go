package models

import (
	"time"
)

type Entry struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	CategoryID uint
	Category   Category
	Rating     int `gorm:"check:rating >= 1 AND rating <= 10"`
	//Notes       string
	ImagePath *string
	Tags      []Tag `gorm:"many2many:entry_tags;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MainCategory struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type Category struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"unique;not null"`
	MainCategoryID uint
	MainCategory   MainCategory
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
