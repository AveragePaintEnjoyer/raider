package web

import (
	"raider/internal/db"
	"raider/internal/models"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", Dashboard)
	app.Get("/categories", Categories)
	app.Get("/sidebar", Sidebar)
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
	return c.Render("dashboard", fiber.Map{})
}

func Sidebar(c *fiber.Ctx) error {
	var mainCats []models.MainCategory
	db.DB.Order("name asc").Find(&mainCats)

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	var entries []models.Entry
	db.DB.Order("name asc").Find(&entries)

	// Category → entries
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

	// MainCategory → categories
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

	var result []MainCategoryWithCategories
	for _, v := range mainMap {
		result = append(result, *v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return c.Render("sidebar", fiber.Map{
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

	var tags []models.Tag
	db.DB.Order("name asc").Find(&tags)

	return c.Render("entry_add", fiber.Map{
		"CategoryID": categoryID,
		"Categories": categories,
		"Tags":       tags,
	})
}

func EntryAdd(c *fiber.Ctx) error {
	name := c.FormValue("name")
	rating, _ := strconv.Atoi(c.FormValue("rating"))
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))

	tagIDs := c.FormValue("tags")
	tags := []models.Tag{}
	for _, s := range strings.Split(tagIDs, ",") {
		if id, err := strconv.Atoi(s); err == nil {
			var tag models.Tag
			db.DB.First(&tag, id)
			tags = append(tags, tag)
		}
	}

	entry := models.Entry{
		Name:       name,
		Rating:     rating,
		CategoryID: uint(categoryID),
		Tags:       tags,
	}

	db.DB.Create(&entry)

	// Return detail view of the newly created entry in right panel
	c.Set("HX-Trigger", "refreshSidebar")
	return c.Render("entry_detail", fiber.Map{
		"Entry": entry,
	}, "")
}

func EntryEdit(c *fiber.Ctx) error {
	id := c.Params("id")

	var entry models.Entry
	if err := db.DB.
		Preload("Category").
		Preload("Tags").
		First(&entry, id).Error; err != nil {
		return c.Status(404).SendString("Entry not found")
	}

	var categories []models.Category
	db.DB.Order("name asc").Find(&categories)

	var tags []models.Tag
	db.DB.Order("name asc").Find(&tags)

	return c.Render("entry_edit", fiber.Map{
		"Entry":      entry,
		"Categories": categories,
		"Tags":       tags,
	})
}

func EntryDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	db.DB.Delete(&models.Entry{}, id)
	c.Set("HX-Trigger", "refreshSidebar")
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

	entry := models.Entry{}
	db.DB.Preload("Tags").First(&entry, id)

	entry.Name = c.FormValue("name")
	entry.Rating = rating
	entry.CategoryID = uint(categoryID)

	tagIDs := c.FormValue("tags")
	tags := []models.Tag{}
	for _, s := range strings.Split(tagIDs, ",") {
		if tid, err := strconv.Atoi(s); err == nil {
			var tag models.Tag
			db.DB.First(&tag, tid)
			tags = append(tags, tag)
		}
	}
	db.DB.Model(&entry).Association("Tags").Replace(tags)

	db.DB.Save(&entry)

	c.Set("HX-Trigger", "refreshSidebar")
	return c.Render("entry_detail", fiber.Map{
		"Entry": entry,
	}, "")
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
