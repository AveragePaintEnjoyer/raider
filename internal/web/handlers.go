package web

import (
	"raider/internal/db"
	"raider/internal/models"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", Dashboard)
	app.Get("/categories", Categories)
	app.Get("/entries/new", EntryNew)
	app.Post("/entries/add", EntryAdd)
	app.Get("/entries/:id", EntryView)
	app.Get("/entries/:id/edit", EntryEdit)
	app.Post("/entries/:id/edit", EntryEditSubmit)
	app.Post("/entries/:id/delete", EntryDelete)
	app.Post("/main-categories/add", MainCategoriesAdd)
	app.Post("/main-categories/delete/:id", MainCategoriesDelete)
	app.Post("/categories/add", CategoriesAdd)
	app.Post("/categories/delete/:id", CategoriesDelete)
	app.Get("/tags", Tags)
	app.Post("/tags/add", TagsAdd)
	app.Post("/tags/delete/:id", TagsDelete)
}

func Dashboard(c *fiber.Ctx) error {
	var mainCats []models.MainCategory
	db.DB.Order("name asc").Find(&mainCats)

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	var entries []models.Entry
	db.DB.Order("name asc").Find(&entries)

	// Map categories → entries
	catMap := make(map[uint]*CategoryWithEntries)
	for _, cat := range categories {
		catMap[cat.ID] = &CategoryWithEntries{
			ID:      cat.ID,
			Name:    cat.Name,
			Entries: []models.Entry{},
		}
	}

	for _, e := range entries {
		if cat, ok := catMap[e.CategoryID]; ok {
			cat.Entries = append(cat.Entries, e)
		}
	}

	// Map main categories → categories
	mainMap := make(map[uint]*MainCategoryWithCategories)
	for _, mc := range mainCats {
		mainMap[mc.ID] = &MainCategoryWithCategories{
			ID:         mc.ID,
			Name:       mc.Name,
			Categories: []CategoryWithEntries{},
		}
	}

	for _, cat := range categories {
		if mc, ok := mainMap[cat.MainCategoryID]; ok {
			if cwe, ok := catMap[cat.ID]; ok {
				mc.Categories = append(mc.Categories, *cwe)
			}
		}
	}

	// Convert to slice
	var result []MainCategoryWithCategories
	for _, v := range mainMap {
		result = append(result, *v)
	}

	// Sort main categories alphabetically
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return c.Render("dashboard", fiber.Map{
		"MainCategories": result,
	})
}

func EntryView(c *fiber.Ctx) error {
	id := c.Params("id")

	var entry models.Entry
	if err := db.DB.
		Preload("Category").
		Preload("Category.MainCategory").
		Preload("Tags").
		First(&entry, id).Error; err != nil {

		return c.Status(404).SendString("Entry not found")
	}

	return c.Render("entry_detail", fiber.Map{
		"Entry": entry,
	}, "")
}

func EntryNew(c *fiber.Ctx) error {
	categoryIDStr := c.Query("category_id")
	var categoryID uint
	if id, err := strconv.Atoi(categoryIDStr); err == nil {
		categoryID = uint(id)
	}

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	return c.Render("entry_add", fiber.Map{
		"CategoryID": categoryID,
		"Categories": categories,
	}, "")
}

func EntryAdd(c *fiber.Ctx) error {
	name := c.FormValue("name")
	rating, _ := strconv.Atoi(c.FormValue("rating"))
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))

	entry := models.Entry{
		Name:       name,
		Rating:     rating,
		CategoryID: uint(categoryID),
	}

	db.DB.Create(&entry)

	// Return detail view of the newly created entry in right panel
	return c.Render("entry_detail", fiber.Map{
		"Entry": entry,
	}, "")
}

func EntryEdit(c *fiber.Ctx) error {
	id := c.Params("id")

	var entry models.Entry
	if err := db.DB.
		Preload("Category").
		First(&entry, id).Error; err != nil {
		return c.Status(404).SendString("Entry not found")
	}

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	return c.Render("entry_edit", fiber.Map{
		"Entry":      entry,
		"Categories": categories,
	}, "")
}

func EntryDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	db.DB.Delete(&models.Entry{}, id)
	return c.SendString(`
		<div class="table-box">
			<p><em>Entry deleted.</em></p>
		</div>
	`)
}

func EntryEditSubmit(c *fiber.Ctx) error {
	id := c.Params("id")

	rating, _ := strconv.Atoi(c.FormValue("rating"))
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))

	db.DB.Model(&models.Entry{}).
		Where("id = ?", id).
		Updates(models.Entry{
			Name:       c.FormValue("name"),
			Rating:     rating,
			CategoryID: uint(categoryID),
		})

	// Reload updated entry view
	return c.Redirect("/entries/" + id)
}

func Categories(c *fiber.Ctx) error {
	var mainCategories []models.MainCategory
	db.DB.Order("name asc").Find(&mainCategories)

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	return c.Render("categories", fiber.Map{
		"MainCategories": mainCategories,
		"Categories":     categories,
	})
}

func MainCategoriesAdd(c *fiber.Ctx) error {
	mc := models.MainCategory{
		Name: c.FormValue("name"),
	}
	db.DB.Create(&mc)
	return c.Redirect("/categories")
}

func MainCategoriesDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	db.DB.Delete(&models.MainCategory{}, id)
	return c.Redirect("/categories")
}

func CategoriesAdd(c *fiber.Ctx) error {
	mainID, _ := strconv.Atoi(c.FormValue("main_category_id"))

	cat := models.Category{
		Name:           c.FormValue("name"),
		MainCategoryID: uint(mainID),
	}
	db.DB.Create(&cat)

	return c.Redirect("/categories")
}

func CategoriesDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	db.DB.Delete(&models.Category{}, id)
	return c.Redirect("/categories")
}

func Tags(c *fiber.Ctx) error {
	var tags []models.Tag
	db.DB.Find(&tags)

	return c.Render("tags", fiber.Map{
		"Tags": tags,
	})
}

func TagsAdd(c *fiber.Ctx) error {
	s := models.Tag{
		Name: c.FormValue("name"),
	}
	db.DB.Create(&s)
	return c.Redirect("/tags")
}

func TagsDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	db.DB.Delete(&models.Tag{}, id)
	return c.Redirect("/tags")
}
