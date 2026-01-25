package web

import "raider/internal/models"

type CategoryWithEntries struct {
	ID      uint
	Name    string
	Entries []models.Entry
}

type MainCategoryWithCategories struct {
	ID         uint
	Name       string
	Categories []CategoryWithEntries
}
