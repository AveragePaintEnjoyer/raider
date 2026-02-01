package web

import (
	"fmt"
	"os"
	"path/filepath"
)

func EntryImagePath(entryID uint, ext string) (fsPath string, webPath string) {
	if ext == "" {
		ext = ".jpg"
	}

	staticRoot := os.Getenv("STATIC_PATH")

	fsPath = filepath.Join(
		staticRoot,
		"uploads",
		"entries",
		fmt.Sprintf("%d%s", entryID, ext),
	)

	webPath = fmt.Sprintf(
		"/static/uploads/entries/%d%s",
		entryID,
		ext,
	)

	return
}
